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
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/model"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/gin-gonic/gin"
)

// setupTestRouter создает тестовый роутер без middleware логирования
func setupTestRouter(cfg *config.ConfigStruct) (*gin.Engine, *service.URLShortnerService, *auth.AuthService) {
	gin.SetMode(gin.TestMode)
	router := gin.New() // Используем New() вместо Default() чтобы избежать логгера

	repo := repository.NewMemoryRepository()
	svc := service.NewURLShortnerService(repo, cfg)
	authService := auth.NewAuthService(cfg.AuthSecretKey)

	return router, svc, authService
}

// Example_shortenURL демонстрирует создание короткого URL через текстовый запрос.
// Эндпоинт: POST /
func Example_shortenURL() {
	// Настройка тестового сервера
	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	router, svc, authService := setupTestRouter(cfg)

	// Регистрируем только нужный хендлер напрямую
	router.POST("/", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		userID := "user123"

		shortURL, _ := svc.CreateShortURL(string(body), userID)
		c.String(http.StatusCreated, cfg.ShortAddress+"/"+shortURL)
	})

	// Создание тестового запроса
	body := "https://practicum.yandex.ru"
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "text/plain")

	// Добавляем токен аутентификации
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
	fmt.Println("Response contains base URL:", len(string(responseBody)) > len(cfg.ShortAddress))
	// Output:
	// Status: 201
	// Response contains base URL: true
}

// Example_shortenURLJSON демонстрирует создание короткого URL через JSON запрос.
// Эндпоинт: POST /api/shorten
func Example_shortenURLJSON() {
	// Настройка тестового сервера
	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	router, svc, authService := setupTestRouter(cfg)

	// Регистрируем хендлер
	router.POST("/api/shorten", func(c *gin.Context) {
		var request model.Request
		json.NewDecoder(c.Request.Body).Decode(&request)

		userID := "user123"
		shortURL, _ := svc.CreateShortURL(request.URL, userID)

		response := model.Response{
			Result: cfg.ShortAddress + "/" + shortURL,
		}
		c.JSON(http.StatusCreated, response)
	})

	// Создание JSON запроса
	requestBody := model.Request{
		URL: "https://practicum.yandex.ru",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Добавляем токен аутентификации
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
	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	router, svc, authService := setupTestRouter(cfg)

	// Регистрируем хендлер
	router.POST("/api/shorten/batch", func(c *gin.Context) {
		var requests []model.BatchRequest
		json.NewDecoder(c.Request.Body).Decode(&requests)

		urls := make([]string, len(requests))
		correlationMap := make(map[string]string)
		for i, request := range requests {
			urls[i] = request.OriginalURL
			correlationMap[request.OriginalURL] = request.CorrelationID
		}

		userID := "user123"
		shortURLsMap, _ := svc.CreateShortURLsBatch(urls, userID)

		responses := make([]model.BatchResponse, 0, len(requests))
		for originalURL, shortURL := range shortURLsMap {
			if correlationID, exists := correlationMap[originalURL]; exists {
				responses = append(responses, model.BatchResponse{
					CorrelationID: correlationID,
					ShortURL:      cfg.ShortAddress + "/" + shortURL,
				})
			}
		}

		c.JSON(http.StatusCreated, responses)
	})

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
	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	router, svc, authService := setupTestRouter(cfg)

	// Сначала создаем короткий URL
	userID := "user123"
	shortURL, _ := svc.CreateShortURL("https://practicum.yandex.ru", userID)

	// Регистрируем хендлер
	router.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		fullURL, err := svc.GetFullURL(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, fullURL)
	})

	// Теперь получаем оригинальный URL
	req := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)
	cookie := authService.CreateSignedCookie(userID)
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
	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	router, svc, authService := setupTestRouter(cfg)

	// Создаем несколько URL для пользователя
	userID := "user123"
	svc.CreateShortURL("https://practicum.yandex.ru", userID)
	svc.CreateShortURL("https://yandex.ru", userID)

	// Регистрируем хендлер
	router.GET("/api/user/urls", func(c *gin.Context) {
		userURLs, err := svc.GetUserURLs(userID)
		if err != nil || len(userURLs) == 0 {
			c.Status(http.StatusNoContent)
			return
		}

		response := make([]model.UserURL, len(userURLs))
		for i, urlData := range userURLs {
			response[i] = model.UserURL{
				ShortURL:    cfg.ShortAddress + "/" + urlData["short_url"],
				OriginalURL: urlData["original_url"],
			}
		}
		c.JSON(http.StatusOK, response)
	})

	// Получаем все URL пользователя
	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	cookie := authService.CreateSignedCookie(userID)
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
	cfg := &config.ConfigStruct{
		ShortAddress:  "http://localhost:8080",
		AuthSecretKey: "test-secret-key",
	}

	router, svc, authService := setupTestRouter(cfg)

	// Создаем URL для удаления
	userID := "user123"
	shortURL1, _ := svc.CreateShortURL("https://practicum.yandex.ru", userID)
	shortURL2, _ := svc.CreateShortURL("https://yandex.ru", userID)

	// Регистрируем хендлер
	router.DELETE("/api/user/urls", func(c *gin.Context) {
		var deleteRequest []string
		json.NewDecoder(c.Request.Body).Decode(&deleteRequest)

		// Запускаем асинхронное удаление
		go func() {
			svc.DeleteURLsBatch(deleteRequest, userID)
		}()

		c.Status(http.StatusAccepted)
	})

	// Удаляем URL
	deleteRequest := []string{shortURL1, shortURL2}
	jsonBody, _ := json.Marshal(deleteRequest)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	cookie := authService.CreateSignedCookie(userID)
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
	cfg := &config.ConfigStruct{
		ShortAddress: "http://localhost:8080",
		AddressDB:    "", // Пустая строка для тестирования без БД
	}

	router, svc, _ := setupTestRouter(cfg)

	// Регистрируем хендлер
	router.GET("/ping", func(c *gin.Context) {
		if err := svc.PingPostgreSQL(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

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
