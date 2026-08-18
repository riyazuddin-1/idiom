package config

import "os"

var (
	PostgresDSN  = os.Getenv("POSTGRES_DSN")
	MongodbUrl   = os.Getenv("MONGODB_URL")
	JWTPublicKey = os.Getenv("JWT_PUBLIC_KEY")
	AuthBaseURL  = os.Getenv("AUTH_BASE_URL")
	DistDir      = getEnv("DIST_DIR", "app/dist")
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
