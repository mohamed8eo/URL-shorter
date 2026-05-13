package model

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateShortURL() (string, error) {
	randomBytes := make([]byte, 6)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
