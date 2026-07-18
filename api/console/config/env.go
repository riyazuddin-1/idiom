package config

import "os"

var (
	PostgresDSN  = os.Getenv("POSTGRES_DSN")
	MongodbUrl   = os.Getenv("MONGODB_URL")
	JWTPublicKey = os.Getenv("JWT_PUBLIC_KEY")
	AuthBaseURL  = os.Getenv("AUTH_BASE_URL")
)
