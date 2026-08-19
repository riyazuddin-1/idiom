package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

const defaultMailingURL = "https://catalog-worker.vercel.app/api/v1/jobs"

type mailingJob struct {
	Event      string            `json:"event"`
	Properties mailingProperties `json:"properties"`
}

type mailingProperties struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// HTTPMailSender sends emails through the external mailing worker.
type HTTPMailSender struct {
	client *http.Client
	url    string
}

// NewHTTPMailSender creates an email sender using the mailing worker.
func NewHTTPMailSender() *HTTPMailSender {
	return &HTTPMailSender{
		client: http.DefaultClient,
		url:    defaultMailingURL,
	}
}

// Send queues one email with the external mailing worker.
func (s *HTTPMailSender) Send(ctx context.Context, message Message) error {
	body := message.TextBody
	if message.HTMLBody != "" {
		body = message.HTMLBody
	}

	payload, err := json.Marshal(mailingJob{
		Event: "email.send",
		Properties: mailingProperties{
			To:      message.To,
			Subject: message.Subject,
			Body:    body,
		},
	})
	if err != nil {
		return fmt.Errorf("encode email job: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create email request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("send email job: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("email worker returned HTTP status %s", response.Status)
	}

	return nil
}
