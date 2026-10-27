package jwt

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math/big"
)

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func PublicKey(privateKeyPEM string) (string, error) {
	privateKey, err := parseRSAPrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}

	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})), nil
}

// keyID returns a stable identifier derived from the public key.
// The identifier remains the same across service restarts.
func keyID(publicKey interface{}) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(publicKeyBytes)
	return base64.RawURLEncoding.EncodeToString(hash[:]), nil
}

func publicJWK(publicKeyPEM string) (JWK, error) {
	publicKey, err := parseRSAPublicKey(publicKeyPEM)
	if err != nil {
		return JWK{}, err
	}

	kid, err := keyID(publicKey)
	if err != nil {
		return JWK{}, err
	}

	return JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: kid,
		N: base64.RawURLEncoding.EncodeToString(
			publicKey.N.Bytes(),
		),
		E: base64.RawURLEncoding.EncodeToString(
			big.NewInt(int64(publicKey.E)).Bytes(),
		),
	}, nil
}

func (j *JWTSettings) PublicJWKS() (JWKS, error) {
	key, err := publicJWK(j.PublicKey)
	if err != nil {
		return JWKS{}, err
	}

	return JWKS{Keys: []JWK{key}}, nil
}
