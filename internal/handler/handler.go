package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/audit"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/auth"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/middleware"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/model"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// Handler представляет обработчик HTTP-запросов для сервиса сокращения URL.
// Структура содержит все необходимые зависимости для обработки запросов:
// сервис бизнес-логики, конфигурацию, сервисы аутентификации и аудита.
type Handler struct {
	Service       *service.URLShortnerService
	Configuration *config.ConfigStruct
	AuthService   *auth.AuthService
	AuditService  *audit.AuditService
}

// NewHandler создает новый экземпляр Handler и регистрирует все маршруты приложения.
// Функция настраивает middleware (сжатие, логирование, аутентификация) и регистрирует
// все эндпоинты API для работы с сокращенными URL.
//
// Параметры:
//   - ginEngine: экземпляр Gin engine для регистрации маршрутов
//   - service: сервис бизнес-логики для работы с URL
//   - configuration: конфигурация приложения
//
// Регистрируемые маршруты:
//   - POST /api/shorten - создание короткого URL (JSON)
//   - POST /api/shorten/batch - пакетное создание коротких URL
//   - POST / - создание короткого URL (text/plain)
//   - GET /:id - получение оригинального URL по короткому
//   - GET /ping - проверка доступности базы данных
//   - GET /api/user/urls - получение всех URL пользователя
//   - DELETE /api/user/urls - удаление URL пользователя
func NewHandler(
	ginEngine *gin.Engine,
	service *service.URLShortnerService,
	configuration *config.ConfigStruct,
) {
	// Создаем сервис аутентификации
	authService := auth.NewAuthService(configuration.AuthSecretKey)

	// Создаем сервис аудита
	auditService := audit.NewAuditService()

	// Подключаем наблюдателей аудита
	if configuration.AuditFile != "" {
		fileObserver := audit.NewFileObserver(configuration.AuditFile)
		auditService.Attach(fileObserver)
		log.Printf("Audit file observer enabled: %s", configuration.AuditFile)
	}

	if configuration.AuditURL != "" {
		remoteObserver := audit.NewRemoteObserver(configuration.AuditURL)
		auditService.Attach(remoteObserver)
		log.Printf("Audit remote observer enabled: %s", configuration.AuditURL)
	}

	handler := &Handler{
		Service:       service,
		Configuration: configuration,
		AuthService:   authService,
		AuditService:  auditService,
	}

	// Добавляем middleware перед регистрацией маршрутов
	ginEngine.Use(middleware.GzipMiddleware())
	ginEngine.Use(middleware.LoggingMiddleware())
	ginEngine.Use(middleware.AuthMiddleware(authService))

	// Регистрируем маршруты
	ginEngine.POST("/api/shorten", handler.SendJSONURL)
	ginEngine.POST("/api/shorten/batch", handler.SendJSONURLBatch)
	ginEngine.POST("/", handler.SendURL)
	ginEngine.GET("/:id", handler.GetURL)
	ginEngine.GET("/ping", handler.Ping)

	ginEngine.GET("/api/user/urls", handler.GetUserURLs)
	ginEngine.DELETE("/api/user/urls", handler.DeleteUserURLs)
}

// handleServiceError обрабатывает ошибки сервиса и отправляет соответствующий текстовый ответ
func (h *Handler) handleServiceError(c *gin.Context, err error, shortURL string) {
	if errors.Is(err, repository.ErrRowExists) {
		c.String(http.StatusConflict, h.Configuration.ShortAddress+"/"+shortURL)
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.Abort()
}

// handleServiceErrorJSON обрабатывает ошибки сервиса и отправляет соответствующий JSON ответ
func (h *Handler) handleServiceErrorJSON(c *gin.Context, err error, shortURL string) {
	if errors.Is(err, repository.ErrRowExists) {
		var response model.Response
		response.Result = h.Configuration.ShortAddress + "/" + shortURL
		c.JSON(http.StatusConflict, response)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.Abort()
}

// handleGenericErrorJSON обрабатывает общие ошибки и отправляет JSON ответ
func (h *Handler) handleGenericErrorJSON(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}

// handleGenericErrorText обрабатывает общие ошибки и отправляет текстовый ответ
func (h *Handler) handleGenericErrorText(c *gin.Context, statusCode int, message string) {
	c.String(statusCode, message)
	c.Abort()
}

// SendURL обрабатывает POST-запрос для создания короткого URL из текста.
// Эндпоинт принимает оригинальный URL в теле запроса в формате text/plain
// и возвращает сокращенный URL в формате text/plain.
//
// Метод: POST
// Путь: /
// Content-Type: text/plain
// Требуется аутентификация: да (через cookie)
//
// Пример запроса:
//
//	POST / HTTP/1.1
//	Content-Type: text/plain
//	Cookie: auth_token=...
//
//	https://example.com/very-long-url
//
// Пример успешного ответа (201 Created):
//
//	http://localhost:8080/abc123
//
// Пример ответа при конфликте (409 Conflict):
//
//	http://localhost:8080/abc123
//
// Возможные ошибки:
//   - 400 Bad Request: неверный Content-Type или ошибка чтения тела запроса
//   - 409 Conflict: URL уже существует в базе
//   - 500 Internal Server Error: внутренняя ошибка сервера
func (h *Handler) SendURL(c *gin.Context) {

	// Проверка Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "text/plain") {
		h.handleGenericErrorText(c, http.StatusBadRequest, "Invalid ContentType, text/plain only")
		return
	}

	// Чтение тела запроса
	body, err := c.GetRawData()
	if err != nil {
		h.handleGenericErrorText(c, http.StatusBadRequest, "Error reading request body")
		return
	}

	// Получаем userID из контекста
	userID, _ := c.Get(middleware.UserIDKey)
	userIDStr := userID.(string)

	// Создание короткой ссылки
	shortURL, err := h.Service.CreateShortURL(string(body), userIDStr)
	if err != nil {
		h.handleServiceError(c, err, shortURL)
		return
	}

	// Логируем событие аудита
	h.AuditService.LogEvent(audit.ActionShorten, userIDStr, string(body))

	// Response: текст с полным URL
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, h.Configuration.ShortAddress+"/"+shortURL)
}

// SendJSONURL обрабатывает POST-запрос для создания короткого URL из JSON.
// Эндпоинт принимает JSON с оригинальным URL и возвращает JSON с сокращенным URL.
//
// Метод: POST
// Путь: /api/shorten
// Content-Type: application/json
// Требуется аутентификация: да (через cookie)
//
// Пример запроса:
//
//	POST /api/shorten HTTP/1.1
//	Content-Type: application/json
//	Cookie: auth_token=...
//
//	{
//	  "url": "https://example.com/very-long-url"
//	}
//
// Пример успешного ответа (201 Created):
//
//	{
//	  "result": "http://localhost:8080/abc123"
//	}
//
// Пример ответа при конфликте (409 Conflict):
//
//	{
//	  "result": "http://localhost:8080/abc123"
//	}
//
// Возможные ошибки:
//   - 400 Bad Request: неверный Content-Type, невалидный JSON или URL
//   - 409 Conflict: URL уже существует в базе
//   - 500 Internal Server Error: внутренняя ошибка сервера
func (h *Handler) SendJSONURL(c *gin.Context) {

	// Проверка Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		h.handleGenericErrorText(c, http.StatusBadRequest, "Invalid ContentType, application/json only")
		return
	}

	// Получаем данные из body
	var request model.Request
	// Странный тест на обязательный юз encoding/json, чисто ради галочки...
	// if err := c.ShouldBindJSON(&request); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	c.Abort()
	// 	return
	// }
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация реквеста
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем userID из контекста
	userID, _ := c.Get(middleware.UserIDKey)
	userIDStr := userID.(string)

	// Создание короткой ссылки
	shortURL, err := h.Service.CreateShortURL(request.URL, userIDStr)
	if err != nil {
		h.handleServiceErrorJSON(c, err, shortURL)
		return
	}

	// Логируем событие аудита
	h.AuditService.LogEvent(audit.ActionShorten, userIDStr, request.URL)

	// Response: JSON с полным URL
	var response model.Response
	response.Result = h.Configuration.ShortAddress + "/" + shortURL

	// Пишем ответ
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
}

// SendJSONURLBatch обрабатывает POST-запрос для пакетного создания коротких URL.
// Эндпоинт принимает массив JSON-объектов с оригинальными URL и correlation_id,
// и возвращает массив с короткими URL и соответствующими correlation_id.
//
// Метод: POST
// Путь: /api/shorten/batch
// Content-Type: application/json
// Требуется аутентификация: да (через cookie)
//
// Пример запроса:
//
//	POST /api/shorten/batch HTTP/1.1
//	Content-Type: application/json
//	Cookie: auth_token=...
//
//	[
//	  {
//	    "correlation_id": "id1",
//	    "original_url": "https://example.com/url1"
//	  },
//	  {
//	    "correlation_id": "id2",
//	    "original_url": "https://example.com/url2"
//	  }
//	]
//
// Пример успешного ответа (201 Created):
//
//	[
//	  {
//	    "correlation_id": "id1",
//	    "short_url": "http://localhost:8080/abc123"
//	  },
//	  {
//	    "correlation_id": "id2",
//	    "short_url": "http://localhost:8080/def456"
//	  }
//	]
//
// Возможные ошибки:
//   - 400 Bad Request: неверный Content-Type, невалидный JSON, пустой массив или невалидные URL
//   - 500 Internal Server Error: внутренняя ошибка сервера
func (h *Handler) SendJSONURLBatch(c *gin.Context) {

	// Проверка Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, "Invalid ContentType, application/json only")
		return
	}

	// Получаем данные из body
	var requests []model.BatchRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&requests); err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Проверяем, что пакет не пустой
	if len(requests) == 0 {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, "Empty batch not allowed")
		return
	}

	// Валидация всех реквестов
	validate := validator.New()
	for _, request := range requests {
		if err := validate.Struct(request); err != nil {
			h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Извлекаем URL из реквестов
	urls := make([]string, len(requests))
	correlationMap := make(map[string]string) // originalURL -> correlationID
	for i, request := range requests {
		urls[i] = request.OriginalURL
		correlationMap[request.OriginalURL] = request.CorrelationID
	}

	// Получаем userID из контекста
	userID, _ := c.Get(middleware.UserIDKey)
	userIDStr := userID.(string)

	// Создание коротких ссылок пакетом
	shortURLsMap, err := h.Service.CreateShortURLsBatch(urls, userIDStr)
	if err != nil {
		h.handleGenericErrorJSON(c, http.StatusInternalServerError, "Error creating short URL")
		return
	}

	// Формируем ответ
	responses := make([]model.BatchResponse, 0, len(requests))
	for originalURL, shortURL := range shortURLsMap {
		if correlationID, exists := correlationMap[originalURL]; exists {
			responses = append(responses, model.BatchResponse{
				CorrelationID: correlationID,
				ShortURL:      h.Configuration.ShortAddress + "/" + shortURL,
			})
		}
	}

	// Пишем ответ
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, responses)
}

// GetURL обрабатывает GET-запрос для получения оригинального URL по короткому.
// Эндпоинт выполняет редирект на оригинальный URL, если короткий URL существует.
//
// Метод: GET
// Путь: /:id
// Требуется аутентификация: нет
//
// Пример запроса:
//
//	GET /abc123 HTTP/1.1
//
// Пример успешного ответа (307 Temporary Redirect):
//
//	HTTP/1.1 307 Temporary Redirect
//	Location: https://example.com/original-url
//
// Возможные ошибки:
//   - 400 Bad Request: короткий URL не указан
//   - 404 Not Found: короткий URL не найден
//   - 410 Gone: URL был удален пользователем
func (h *Handler) GetURL(c *gin.Context) {
	// Получаем параметр из URL: /:id
	shortURL := c.Param("id")

	// Проверяем, удален ли URL
	if deleted, err := h.Service.IsDeleted(shortURL); err == nil && deleted {
		c.Status(http.StatusGone)
		return
	}

	// Ищем полную ссылку
	fullURL, err := h.Service.GetFullURL(shortURL)
	if err != nil {
		h.handleGenericErrorText(c, http.StatusBadRequest, "URL not found")
		return
	}

	// Получаем userID из контекста (может быть пустым для анонимных пользователей)
	userID, exists := c.Get(middleware.UserIDKey)
	userIDStr := ""
	if exists {
		userIDStr = userID.(string)
	}

	// Логируем событие аудита
	h.AuditService.LogEvent(audit.ActionFollow, userIDStr, fullURL)

	// Редирект (307)
	c.Redirect(http.StatusTemporaryRedirect, fullURL)
}

// Ping проверяет доступность базы данных PostgreSQL.
// Эндпоинт используется для healthcheck и мониторинга состояния приложения.
//
// Метод: GET
// Путь: /ping
// Требуется аутентификация: нет
//
// Пример запроса:
//
//	GET /ping HTTP/1.1
//
// Пример успешного ответа (200 OK):
//
//	{
//	  "status": "OK"
//	}
//
// Возможные ошибки:
//   - 500 Internal Server Error: база данных недоступна
func (h *Handler) Ping(c *gin.Context) {
	if err := h.Service.PingPostgreSQL(); err != nil {
		h.handleGenericErrorJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// GetUserURLs возвращает все URL, созданные аутентифицированным пользователем.
// Эндпоинт предоставляет список всех коротких и оригинальных URL пользователя.
//
// Метод: GET
// Путь: /api/user/urls
// Требуется аутентификация: да (через cookie)
//
// Пример запроса:
//
//	GET /api/user/urls HTTP/1.1
//	Cookie: auth_token=...
//
// Пример успешного ответа (200 OK):
//
//	[
//	  {
//	    "short_url": "http://localhost:8080/abc123",
//	    "original_url": "https://example.com/url1"
//	  },
//	  {
//	    "short_url": "http://localhost:8080/def456",
//	    "original_url": "https://example.com/url2"
//	  }
//	]
//
// Возможные ошибки:
//   - 204 No Content: у пользователя нет созданных URL
//   - 401 Unauthorized: пользователь не аутентифицирован
//   - 500 Internal Server Error: внутренняя ошибка сервера
func (h *Handler) GetUserURLs(c *gin.Context) {
	// Получаем userID из контекста
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		h.handleGenericErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDStr := userID.(string)
	if userIDStr == "" {
		h.handleGenericErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Получаем URL пользователя
	userURLs, err := h.Service.GetUserURLs(userIDStr)
	if err != nil {
		h.handleGenericErrorJSON(c, http.StatusInternalServerError, "Error retrieving user URLs")
		return
	}

	// Если у пользователя нет URL
	if len(userURLs) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// Формируем ответ с полными URL
	response := make([]model.UserURL, len(userURLs))
	for i, urlData := range userURLs {
		response[i] = model.UserURL{
			ShortURL:    h.Configuration.ShortAddress + "/" + urlData["short_url"],
			OriginalURL: urlData["original_url"],
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

// DeleteUserURLs помечает URL как удаленные для указанного пользователя.
// Эндпоинт принимает массив коротких URL для удаления (пометки как deleted).
// Удаление происходит асинхронно в фоновом режиме.
//
// Метод: DELETE
// Путь: /api/user/urls
// Content-Type: application/json
// Требуется аутентификация: да (через cookie)
//
// Пример запроса:
//
//	DELETE /api/user/urls HTTP/1.1
//	Content-Type: application/json
//	Cookie: auth_token=...
//
//	[
//	  "abc123",
//	  "def456"
//	]
//
// Пример успешного ответа (202 Accepted):
//
//	HTTP/1.1 202 Accepted
//
// Возможные ошибки:
//   - 400 Bad Request: неверный Content-Type, невалидный JSON или пустой массив
//   - 401 Unauthorized: пользователь не аутентифицирован
func (h *Handler) DeleteUserURLs(c *gin.Context) {
	// Проверка Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, "Invalid ContentType, application/json only")
		return
	}

	// Получаем userID из контекста
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		h.handleGenericErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDStr := userID.(string)
	if userIDStr == "" {
		h.handleGenericErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Получаем данные из body
	var deleteRequest model.DeleteURLsRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&deleteRequest); err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Проверяем, что список не пустой
	if len(deleteRequest) == 0 {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, "Empty list not allowed")
		return
	}

	// Запускаем асинхронное удаление URL с обработкой ошибок
	go func() {
		if err := h.Service.DeleteURLsBatch(deleteRequest, userIDStr); err != nil {
			log.Printf("Error deleting URLs batch for user %s: %v", userIDStr, err)
		}
	}()

	// Возвращаем статус 202 Accepted
	c.Status(http.StatusAccepted)
}
