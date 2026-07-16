package config

import "os"

var (
	PostgresDSN = os.Getenv("POSTGRES_DSN")
	MongodbUrl  = os.Getenv("MONGODB_URL")
)
