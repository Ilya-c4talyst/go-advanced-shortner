package utils

import (
	"regexp"
	"testing"
)

// Проверка длины и формата сгенерированной ссылки
func TestGenerateShortKey(t *testing.T) {
	key := GenerateShortKey()

	if len(key) != 6 {
		t.Fatalf("ожидалась длина 6, получено %d: %s", len(key), key)
	}

	// Ключ содержит только разрешённые символы
	validChars := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validChars.MatchString(key) {
		t.Errorf("ключ содержит недопустимые символы: %s", key)
	}
}
