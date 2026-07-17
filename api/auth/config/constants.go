package config

import "net/url"

const (
	OperationEmailVerification = "email_verification"
	OperationPasswordReset     = "password_reset"

	TokenExpirationInSeconds = 600
	TokenIssuer              = "idiom:auth"

	VerificationPath = "/auth/v1/verify"
)

func PasswordResetPath(projectID string) string {
	return "/" + url.PathEscape(projectID) + "/password-reset"
}
