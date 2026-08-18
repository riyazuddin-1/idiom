package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashString(value string) (string, error) {
	bytes := sha256.Sum256([]byte(value))

	// 2. Convert the byte array to a hexadecimal string
	return hex.EncodeToString(bytes[:]), nil
}

func HashJSON(value interface{}) (string, error) {
	plaintext, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal json for hashing: %w", err)
	}
	// Calculate SHA256 checksum
	hash := sha256.Sum256(plaintext)

	// Return hex encoded string (you can also use base64 if it fits your project better)
	return hex.EncodeToString(hash[:]), nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
