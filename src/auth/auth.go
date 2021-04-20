package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func GenerateSecureKey(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("Error generating key")
	}
	return hex.EncodeToString(b), nil
}
