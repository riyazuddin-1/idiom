package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenerateID(prefix string, length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return prefix + hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}

	return prefix + hex.EncodeToString(bytes)
}

func RandomString(length int) (string, error) {
	// Hex encoding doubles the byte length, so we divide by 2
	byteLength := (length + 1) / 2
	bytes := make([]byte, byteLength)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to string and truncate to the exact requested length
	return hex.EncodeToString(bytes)[:length], nil
}
