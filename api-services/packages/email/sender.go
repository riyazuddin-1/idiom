package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// Message contains the information required to send an email.
//
// If HTMLBody is provided, it is used as the email body.
// Otherwise, TextBody is used.
type Message struct {
	To       string
	Subject  string
	TextBody string
	HTMLBody string
}

// Sender abstracts email delivery.
//
// This makes the rest of the application independent of the
// underlying email provider or SMTP implementation.
type Sender interface {
	Send(ctx context.Context, message Message) error
}

// SMTPConfig contains the SMTP server configuration.
//
// Security determines how the SMTP connection is secured:
//
//   - "tls"      -> implicit TLS, commonly used on port 465
//   - "starttls" -> SMTP connection first, then STARTTLS,
//     commonly used on port 587
//   - "none"     -> plain SMTP connection (generally not recommended)
//
// Username and Password can be left empty when the SMTP server
// does not require authentication.
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	Security string
}

// SMTPSender sends emails through an SMTP server.
type SMTPSender struct {
	config SMTPConfig
}

// NewSMTPSender creates an SMTP email sender using the supplied
// server configuration.
func NewSMTPSender(config SMTPConfig) *SMTPSender {
	return &SMTPSender{
		config: config,
	}
}

// Send sends one email.
//
// The SMTP configuration is shared regardless of the security mode.
// The only difference is how the connection to the SMTP server is
// established.
func (s *SMTPSender) Send(ctx context.Context, message Message) error {
	// Check whether the caller has already cancelled the operation
	// before starting any network work.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Prefer HTML when it is provided. Otherwise send plain text.
	body := message.TextBody
	contentType := "text/plain; charset=UTF-8"

	if message.HTMLBody != "" {
		body = message.HTMLBody
		contentType = "text/html; charset=UTF-8"
	}

	// Construct the raw MIME email.
	raw := strings.Join([]string{
		"From: " + s.config.From,
		"To: " + message.To,
		"Subject: " + message.Subject,
		"MIME-Version: 1.0",
		"Content-Type: " + contentType,
		"",
		body,
	}, "\r\n")

	addr := net.JoinHostPort(s.config.Host, s.config.Port)

	// SMTP authentication is optional.
	//
	// PlainAuth is appropriate here because the SMTP connection
	// must already be protected by TLS/STARTTLS before credentials
	// are sent.
	var auth smtp.Auth

	if s.config.Username != "" || s.config.Password != "" {
		auth = smtp.PlainAuth(
			"",
			s.config.Username,
			s.config.Password,
			s.config.Host,
		)
	}

	switch strings.ToLower(s.config.Security) {
	case "tls":
		// Implicit TLS: TLS is established immediately when the
		// TCP connection is opened. This is commonly used with 465.
		return s.sendImplicitTLS(ctx, addr, auth, message.To, raw)

	case "starttls":
		// STARTTLS: establish a normal SMTP connection first and
		// upgrade it to TLS using the SMTP STARTTLS command.
		return s.sendSTARTTLS(ctx, addr, auth, message.To, raw)

	case "none":
		// Plain SMTP without TLS.
		//
		// This should generally only be used for trusted internal
		// SMTP servers.
		return s.sendPlain(ctx, addr, auth, message.To, raw)

	default:
		return fmt.Errorf("unsupported SMTP security mode %q", s.config.Security)
	}
}

// sendPlain sends an email over a normal SMTP connection.
//
// This is normally not recommended for internet-facing SMTP because
// credentials and email contents may not be encrypted.
func (s *SMTPSender) sendPlain(
	ctx context.Context,
	addr string,
	auth smtp.Auth,
	to string,
	raw string,
) error {
	client, err := s.newClient(ctx, addr)
	if err != nil {
		return err
	}
	defer client.Close()

	return s.sendMessage(client, auth, to, raw)
}

// sendSTARTTLS connects to the SMTP server normally and then upgrades
// the connection to TLS using STARTTLS.
//
// Port 587 commonly uses this mode.
func (s *SMTPSender) sendSTARTTLS(
	ctx context.Context,
	addr string,
	auth smtp.Auth,
	to string,
	raw string,
) error {
	client, err := s.newClient(ctx, addr)
	if err != nil {
		return err
	}
	defer client.Close()

	// Upgrade the existing SMTP connection to TLS.
	tlsConfig := &tls.Config{
		ServerName: s.config.Host,
	}

	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("start TLS: %w", err)
	}

	return s.sendMessage(client, auth, to, raw)
}

// sendImplicitTLS sends an email over an SMTP connection where TLS
// is established before the SMTP protocol starts.
//
// Port 465 commonly uses implicit TLS.
func (s *SMTPSender) sendImplicitTLS(
	ctx context.Context,
	addr string,
	auth smtp.Auth,
	to string,
	raw string,
) error {
	// Establish the TLS connection immediately.
	tlsConfig := &tls.Config{
		ServerName: s.config.Host,
	}

	conn, err := (&tls.Dialer{
		Config: tlsConfig,
	}).DialContext(ctx, "tcp", addr)

	if err != nil {
		return fmt.Errorf("connect to SMTP server: %w", err)
	}

	// smtp.NewClient takes an already-established connection.
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer client.Close()

	return s.sendMessage(client, auth, to, raw)
}

// newClient establishes a normal TCP connection to the SMTP server.
//
// DialContext is used instead of net.Dial so that the caller's
// context can cancel a connection attempt.
func (s *SMTPSender) newClient(
	ctx context.Context,
	addr string,
) (*smtp.Client, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("connect to SMTP server: %w", err)
	}

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("create SMTP client: %w", err)
	}

	return client, nil
}

// sendMessage performs the SMTP commands required to deliver the
// already-constructed email.
//
// The flow is:
//
//	AUTH
//	  ↓
//	MAIL FROM
//	  ↓
//	RCPT TO
//	  ↓
//	DATA
//	  ↓
//	email contents
//	  ↓
//	QUIT
func (s *SMTPSender) sendMessage(
	client *smtp.Client,
	auth smtp.Auth,
	to string,
	raw string,
) error {
	// Authenticate if credentials were configured.
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication: %w", err)
		}
	}

	// MAIL FROM identifies the sender/envelope address.
	if err := client.Mail(s.config.From); err != nil {
		return fmt.Errorf("SMTP MAIL FROM: %w", err)
	}

	// RCPT TO identifies the actual recipient.
	//
	// This must use message.To, NOT s.config.From.
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT TO: %w", err)
	}

	// DATA begins the email body.
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA: %w", err)
	}

	if _, err := writer.Write([]byte(raw)); err != nil {
		writer.Close()
		return fmt.Errorf("write email: %w", err)
	}

	// Closing DATA tells the SMTP server that the email contents
	// are complete.
	if err := writer.Close(); err != nil {
		return fmt.Errorf("finish email: %w", err)
	}

	// QUIT cleanly closes the SMTP session.
	if err := client.Quit(); err != nil {
		return fmt.Errorf("close SMTP session: %w", err)
	}

	return nil
}
