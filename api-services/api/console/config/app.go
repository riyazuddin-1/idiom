package config

import (
	"idiom-api-services/packages/database/postgres"
	"idiom-api-services/packages/jwt"
)

type AppConfig struct {
	PostgresDB  *postgres.Postgres
	JWTSettings *jwt.JWTSettings
	AuthBaseURL string
	AuthProject string
}
