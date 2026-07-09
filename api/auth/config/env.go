package config

import "os"

var (
	JWTSecret          = os.Getenv("JWT_SECRET")
	AuthBaseURL        = os.Getenv("AUTH_BASE_URL")
	VerificationSecret = os.Getenv("VERIFICATION_SECRET")
	SMTPHost           = os.Getenv("SMTP_HOST")
	SMTPPort           = os.Getenv("SMTP_PORT")
	SMTPUsername       = os.Getenv("SMTP_USERNAME")
	SMTPPassword       = os.Getenv("SMTP_PASSWORD")
	SMTPFrom           = os.Getenv("SMTP_FROM")
)
