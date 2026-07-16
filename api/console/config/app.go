package config

import "idiom-api-services/packages/database/postgres"

type AppConfig struct {
	PostgresDB *postgres.Postgres
}
