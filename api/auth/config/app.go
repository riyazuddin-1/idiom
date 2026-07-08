package config

import "idiom-api-services/packages/jwt"

type AppConfig struct {
	JWTSettings *jwt.JWTSettings
}
