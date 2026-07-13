package config

const (
	ProjectID         = "idiom"
	VerificationPath  = "/auth/v1/verify"
	PasswordResetPath = "/web/password-reset"

	OperationEmailVerification = "email_verification"
	OperationPasswordReset     = "password_reset"

	TokenExpirationInSeconds = 600
	TokenIssuer              = "idiom:auth"
)
