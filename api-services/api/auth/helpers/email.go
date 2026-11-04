package helpers

import (
	"context"
	"idiom-api-services/api/auth/config"
	"idiom-api-services/domains/identities"
	"idiom-api-services/packages/crypto"
	"idiom-api-services/packages/email"
	"net/url"
	"strings"
)

type Scope struct {
	Operation string `json:"operation"`
	ProjectID string `json:"project_id"`
	Email     string `json:"email"`
}

func SendVerificationEmail(ctx context.Context, identity *identities.Identity) error {
	scope, err := crypto.EncryptJSON(Scope{
		Operation: config.OperationEmailVerification,
		ProjectID: identity.ProjectID,
		Email:     identity.Email,
	}, config.VerificationSecret)
	if err != nil {
		return err
	}

	verifyURL := strings.TrimRight(config.AuthBaseURL, "/") + config.VerificationPath + "?scope=" + url.QueryEscape(scope)
	return config.EmailSender.Send(ctx, email.Message{
		To:       identity.Email,
		Subject:  "Verify your email",
		TextBody: "Verify your email by opening this link:\n\n" + verifyURL,
		HTMLBody: `<p>Verify your email by opening this link:</p><p><a href="` + verifyURL + `">Verify email</a></p>`,
	})
}

func SendPasswordResetEmail(ctx context.Context, projectID, emailAddress string) error {
	scope, err := crypto.EncryptJSON(Scope{
		Operation: config.OperationPasswordReset,
		ProjectID: projectID,
		Email:     emailAddress,
	}, config.VerificationSecret)
	if err != nil {
		return err
	}

	resetURL := strings.TrimRight(config.AuthBaseURL, "/") + config.PasswordResetPath(projectID) + "?scope=" + url.QueryEscape(scope)
	return config.EmailSender.Send(ctx, email.Message{
		To:       emailAddress,
		Subject:  "Reset your password",
		TextBody: "Reset your password by opening this link:\n\n" + resetURL,
		HTMLBody: `<p>Reset your password by opening this link:</p><p><a href="` + resetURL + `">Reset password</a></p>`,
	})
}
