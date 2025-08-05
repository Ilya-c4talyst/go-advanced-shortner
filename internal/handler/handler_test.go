package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Сет-ап для тестов
func setupTest() (*gin.Engine, *Handler) {
	db := storage.CreateDB()
	repo := repository.NewShortenerRepository(db)
	service := service.NewURLShortnerService(repo)
	ginEngine := gin.Default()
	h := &Handler{Service: service}

	// Инициализируем роуты
	NewHandler(ginEngine, service)

	return ginEngine, h
}

// Тесты для создания ссылки
func TestSendURLHandler(t *testing.T) {
	mux, _ := setupTest()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Создание ссылки
	t.Run("successful short URL creation", func(t *testing.T) {
		longURL := "https://example.com/very/long/url"
		req, err := http.NewRequest("POST", server.URL+"/", bytes.NewBufferString(longURL))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))

		_, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)
	})

	// Тест на тип контента
	t.Run("invalid content type", func(t *testing.T) {
		req, err := http.NewRequest("POST", server.URL+"/", bytes.NewBufferString("https://test.com"))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// Тест на пустой реквест
	t.Run("empty body", func(t *testing.T) {
		req, err := http.NewRequest("POST", server.URL+"/", bytes.NewBufferString(""))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}

// Тесты на получение полной ссылки
func TestGetURLHandler(t *testing.T) {
	mux, h := setupTest()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Предварительно создаем тестовую короткую ссылку
	longURL := "https://redirect.me"
	shortURL := h.Service.CreateShortURL(longURL)

	// Создаем клиент, который не следует за редиректами автоматически
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Останавливаемся после первого редиректа
		},
	}

	// Тесты редиректов
	// 1)
	t.Run("successful redirect", func(t *testing.T) {
		req, err := http.NewRequest("GET", server.URL+"/"+shortURL, nil)
		assert.NoError(t, err)

		resp, err := client.Do(req) // Используем наш клиент без авто-редиректов
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		assert.Equal(t, longURL, resp.Header.Get("Location"))
	})

	// 2)
	t.Run("not found redirect", func(t *testing.T) {
		req, err := http.NewRequest("GET", server.URL+"/nonexistent", nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
