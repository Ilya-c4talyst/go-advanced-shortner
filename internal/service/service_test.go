package service

import (
	"testing"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestURLShortnerService(t *testing.T) {
	// Инициализируем реальные зависимости
	db := storage.CreateDB()
	repo := repository.NewShortenerRepository(db)
	service := NewURLShortnerService(repo)

	t.Run("Create and get short URL", func(t *testing.T) {
		originalURL := "https://example.com/very/long/url"

		// Создаем короткую ссылку
		shortURL := service.CreateShortURL(originalURL)
		assert.NotEmpty(t, shortURL)
		assert.Len(t, shortURL, 6)

		// Получаем оригинальную ссылку
		fullURL, err := service.GetFullURL(shortURL)
		assert.NoError(t, err)
		assert.Equal(t, originalURL, fullURL)
	})

	t.Run("Get non-existent short URL returns error", func(t *testing.T) {
		nonExistentKey := "nonexist"

		// Пытаемся получить несуществующую ссылку
		_, err := service.GetFullURL(nonExistentKey)
		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
	})

	t.Run("Generate unique short URLs", func(t *testing.T) {
		url1 := "https://first.com"
		url2 := "https://second.com"

		// Генерируем две короткие ссылки
		short1 := service.CreateShortURL(url1)
		short2 := service.CreateShortURL(url2)

		assert.NotEqual(t, short1, short2)

		// Проверяем, что они ведут на разные URL
		full1, _ := service.GetFullURL(short1)
		full2, _ := service.GetFullURL(short2)

		assert.Equal(t, url1, full1)
		assert.Equal(t, url2, full2)
	})

	t.Run("Empty URL handling", func(t *testing.T) {
		emptyURL := ""

		// Не должно паниковать при пустом URL
		shortURL := service.CreateShortURL(emptyURL)
		assert.NotEmpty(t, shortURL)

		fullURL, err := service.GetFullURL(shortURL)
		assert.NoError(t, err)
		assert.Equal(t, emptyURL, fullURL)
	})
}
