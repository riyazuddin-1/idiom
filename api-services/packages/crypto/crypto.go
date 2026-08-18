package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"
)

func GenerateID(prefix string) string {
	timestamp := time.Now().UnixMilli()

	// Converting timestamp to base16(hex)
	hexStr := strconv.FormatInt(timestamp, 16)

	return prefix + hexStr
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
