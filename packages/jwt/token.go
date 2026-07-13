package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTSettings struct {
	PrivateKey          string
	PublicKey           string
	ExpirationInSeconds int64
	Issuer              string
}

type TokenClaims struct {
	Email     string
	ProjectID string
}

func NewJWTSettings(privateKey, publicKey string, expiration int64, issuer string) *JWTSettings {
	return &JWTSettings{
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		ExpirationInSeconds: expiration,
		Issuer:              issuer,
	}
}

func (j *JWTSettings) CreateToken(email, projectId string) (string, error) {
	claims := jwt.MapClaims{
		"email":     email,
		"projectId": projectId,
		"iss":       j.Issuer,
		"exp":       time.Now().Unix() + j.ExpirationInSeconds,
		"iat":       time.Now().Unix(),
		"sub":       email,
	}
	privateKey, err := parseRSAPrivateKey(j.PrivateKey)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTSettings) VerifyToken(tokenString string) (*TokenClaims, error) {
	publicKey, err := parseRSAPublicKey(j.PublicKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims["iss"] != j.Issuer {
		return nil, errors.New("invalid token")
	}

	email, _ := claims["email"].(string)
	projectID, _ := claims["projectId"].(string)

	return &TokenClaims{
		Email:     email,
		ProjectID: projectID,
	}, nil
}

func parseRSAPrivateKey(value string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(value))
	if block == nil {
		return nil, errors.New("invalid private key")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}

	return privateKey, nil
}

func parseRSAPublicKey(value string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(value))
	if block == nil {
		return nil, errors.New("invalid public key")
	}

	if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return key, nil
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}

	return publicKey, nil
}
