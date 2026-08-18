package config

import (
	"idiom-api-services/packages/database/postgres"
	"idiom-api-services/packages/email"
	"idiom-api-services/packages/jwt"
)

type AppConfig struct {
	JWTSettings        *jwt.JWTSettings
	EmailSender        email.Sender
	AuthBaseURL        string
	VerificationSecret string
	PostgresDB         *postgres.Postgres
}
