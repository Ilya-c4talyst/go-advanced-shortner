package examples

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/auth"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/handler"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/model"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/gin-gonic/gin"
)

// Example_shortenURL демонстрирует создание короткого URL через текстовый запрос.
// Эндпоинт: POST /
func Example_shortenURL() {
	// Настройка тестового сервера
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	repo := repository.NewMemoryRepository()
	svc := service.NewURLShortnerService(repo, cfg)
	handler.NewHandler(router, svc, cfg)

	// Создание тестового запроса
	body := "https://practicum.yandex.ru"
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "text/plain")

	// Добавляем токен аутентификации
	authService := auth.NewAuthService(cfg.AuthSecretKey)
	userID := "user123"
	cookie := authService.CreateSignedCookie(userID)
	req.AddCookie(cookie)

	// Выполнение запроса
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Получение ответа
	resp := w.Result()
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Response:", string(responseBody))
	// Output:
	// Status: 201
	// Response: http://localhost:8080/...
}

// Example_shortenURLJSON демонстрирует создание короткого URL через JSON запрос.
// Эндпоинт: POST /api/shorten
func Example_shortenURLJSON() {
	// Настройка тестового сервера
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	repo := repository.NewMemoryRepository()
	svc := service.NewURLShortnerService(repo, cfg)
	handler.NewHandler(router, svc, cfg)

	// Создание JSON запроса
	requestBody := model.Request{
		URL: "https://practicum.yandex.ru",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Добавляем токен аутентификации
	authService := auth.NewAuthService(cfg.AuthSecretKey)
	userID := "user123"
	cookie := authService.CreateSignedCookie(userID)
	req.AddCookie(cookie)

	// Выполнение запроса
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Получение ответа
	resp := w.Result()
	defer resp.Body.Close()

	var response model.Response
	json.NewDecoder(resp.Body).Decode(&response)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Has short URL:", len(response.Result) > 0)
	// Output:
	// Status: 201
	// Has short URL: true
}

// Example_shortenBatch демонстрирует пакетное создание коротких URL.
// Эндпоинт: POST /api/shorten/batch
func Example_shortenBatch() {
	// Настройка тестового сервера
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	repo := repository.NewMemoryRepository()
	svc := service.NewURLShortnerService(repo, cfg)
	handler.NewHandler(router, svc, cfg)

	// Создание пакетного запроса
	batchRequest := []model.BatchRequest{
		{
			CorrelationID: "req1",
			OriginalURL:   "https://practicum.yandex.ru",
		},
		{
			CorrelationID: "req2",
			OriginalURL:   "https://yandex.ru",
		},
	}
	jsonBody, _ := json.Marshal(batchRequest)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Добавляем токен аутентификации
	authService := auth.NewAuthService(cfg.AuthSecretKey)
	userID := "user123"
	cookie := authService.CreateSignedCookie(userID)
	req.AddCookie(cookie)

	// Выполнение запроса
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Получение ответа
	resp := w.Result()
	defer resp.Body.Close()

	var response []model.BatchResponse
	json.NewDecoder(resp.Body).Decode(&response)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Response count:", len(response))
	fmt.Println("First correlation ID:", response[0].CorrelationID)
	// Output:
	// Status: 201
	// Response count: 2
	// First correlation ID: req1
}

// Example_getURL демонстрирует получение оригинального URL по короткому.
// Эндпоинт: GET /:id
func Example_getURL() {
	// Настройка тестового сервера
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	repo := repository.NewMemoryRepository()
	svc := service.NewURLShortnerService(repo, cfg)
	handler.NewHandler(router, svc, cfg)

	// Сначала создаем короткий URL
	authService := auth.NewAuthService(cfg.AuthSecretKey)
	userID := "user123"
	cookie := authService.CreateSignedCookie(userID)

	shortURL, _ := svc.CreateShortURL("https://practicum.yandex.ru", userID)

	// Теперь получаем оригинальный URL
	req := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)
	req.AddCookie(cookie)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Location:", resp.Header.Get("Location"))
	// Output:
	// Status: 307
	// Location: https://practicum.yandex.ru
}

// Example_getUserURLs демонстрирует получение всех URL пользователя.
// Эндпоинт: GET /api/user/urls
func Example_getUserURLs() {
	// Настройка тестового сервера
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	repo := repository.NewMemoryRepository()
	svc := service.NewURLShortnerService(repo, cfg)
	handler.NewHandler(router, svc, cfg)

	// Создаем несколько URL для пользователя
	authService := auth.NewAuthService(cfg.AuthSecretKey)
	userID := "user123"
	cookie := authService.CreateSignedCookie(userID)

	svc.CreateShortURL("https://practicum.yandex.ru", userID)
	svc.CreateShortURL("https://yandex.ru", userID)

	// Получаем все URL пользователя
	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req.AddCookie(cookie)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	var userURLs []model.UserURL
	json.NewDecoder(resp.Body).Decode(&userURLs)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("URL count:", len(userURLs))
	// Output:
	// Status: 200
	// URL count: 2
}

// Example_deleteUserURLs демонстрирует удаление URL пользователя.
// Эндпоинт: DELETE /api/user/urls
func Example_deleteUserURLs() {
	// Настройка тестового сервера
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	repo := repository.NewMemoryRepository()
	svc := service.NewURLShortnerService(repo, cfg)
	handler.NewHandler(router, svc, cfg)

	// Создаем URL для удаления
	authService := auth.NewAuthService(cfg.AuthSecretKey)
	userID := "user123"
	cookie := authService.CreateSignedCookie(userID)

	shortURL1, _ := svc.CreateShortURL("https://practicum.yandex.ru", userID)
	shortURL2, _ := svc.CreateShortURL("https://yandex.ru", userID)

	// Удаляем URL
	deleteRequest := []string{shortURL1, shortURL2}
	jsonBody, _ := json.Marshal(deleteRequest)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println("Status:", resp.StatusCode)
	// Output:
	// Status: 202
}

// Example_ping демонстрирует проверку доступности базы данных.
// Эндпоинт: GET /ping
func Example_ping() {
	// Настройка тестового сервера
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.ConfigStruct{
		ShortAddress: "http://localhost:8080",
		AddressDB:    "", // Пустая строка для тестирования без БД
	}

	repo := repository.NewMemoryRepository()
	svc := service.NewURLShortnerService(repo, cfg)
	handler.NewHandler(router, svc, cfg)

	// Проверка доступности
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	fmt.Println("Status code:", resp.StatusCode)
	// Output будет зависеть от наличия БД
	// В тесте без БД будет ошибка
}
