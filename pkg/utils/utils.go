package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// Генерация короткого ключа (6 символов)
func GenerateShortKey() string {
	b := make([]byte, 4)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:6]
}
