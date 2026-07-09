package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type Message struct {
	To       string
	Subject  string
	TextBody string
	HTMLBody string
}

type Sender interface {
	Send(ctx context.Context, message Message) error
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type SMTPSender struct {
	config SMTPConfig
}

func NewSMTPSender(config SMTPConfig) *SMTPSender {
	config.Host = strings.TrimSpace(config.Host)
	config.Port = strings.TrimSpace(config.Port)
	config.From = strings.TrimSpace(config.From)

	return &SMTPSender{config: config}
}

func (s *SMTPSender) Send(ctx context.Context, message Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	to := strings.TrimSpace(message.To)
	body := message.TextBody
	contentType := "text/plain; charset=UTF-8"
	if strings.TrimSpace(message.HTMLBody) != "" {
		body = message.HTMLBody
		contentType = "text/html; charset=UTF-8"
	}

	raw := strings.Join([]string{
		"From: " + s.config.From,
		"To: " + to,
		"Subject: " + message.Subject,
		"MIME-Version: 1.0",
		"Content-Type: " + contentType,
		"",
		body,
	}, "\r\n")

	addr := s.config.Host + ":" + s.config.Port
	var auth smtp.Auth
	if s.config.Username != "" || s.config.Password != "" {
		auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	}

	if err := smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(raw)); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
