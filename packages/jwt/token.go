package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTSettings struct {
	SecretKey string
}

func NewJWTSettings(secretKey string) *JWTSettings {
	return &JWTSettings{
		SecretKey: secretKey,
	}
}

func (j *JWTSettings) CreateToken(email, projectId string) (string, error) {
	claims := jwt.MapClaims{
		"email":     email,
		"projectId": projectId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
