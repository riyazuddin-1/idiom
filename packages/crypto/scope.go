package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func EncryptJSON(value interface{}, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("encryption secret is required")
	}

	plaintext, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(scopeKey(secret))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func DecryptJSON(scope, secret string, target interface{}) error {
	if secret == "" {
		return errors.New("encryption secret is required")
	}

	payload, err := base64.RawURLEncoding.DecodeString(scope)
	if err != nil {
		return fmt.Errorf("decode scope: %w", err)
	}

	block, err := aes.NewCipher(scopeKey(secret))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(payload) < nonceSize {
		return errors.New("scope is too short")
	}

	plaintext, err := gcm.Open(nil, payload[:nonceSize], payload[nonceSize:], nil)
	if err != nil {
		return fmt.Errorf("decrypt scope: %w", err)
	}

	return json.Unmarshal(plaintext, target)
}

func scopeKey(secret string) []byte {
	sum := sha256.Sum256([]byte(secret))
	return sum[:]
}
