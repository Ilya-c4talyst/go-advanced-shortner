package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/middleware"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/model"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// Handler — структура хендлера
type Handler struct {
	Service       *service.URLShortnerService
	Configuration *config.ConfigStruct
}

// Конструктор для хендлера
func NewHandler(
	ginEngine *gin.Engine,
	service *service.URLShortnerService,
	configuration *config.ConfigStruct,
) {
	handler := &Handler{
		Service:       service,
		Configuration: configuration,
	}

	// Добавляем middleware перед регистрацией маршрутов
	ginEngine.Use(middleware.GzipMiddleware())
	ginEngine.Use(middleware.LoggingMiddleware())

	// Регистрируем маршруты
	ginEngine.POST("/api/shorten", handler.SendJSONURL)
	ginEngine.POST("/", handler.SendURL)
	ginEngine.GET("/:id", handler.GetURL)
	ginEngine.GET("/ping", handler.Ping)
}

// Обработка POST запроса: сокращение URL
func (h *Handler) SendURL(c *gin.Context) {

	// Проверка Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "text/plain") {
		c.String(http.StatusBadRequest, "Invalid ContentType, text/plain only")
		c.Abort()
		return
	}

	// Чтение тела запроса
	body, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, "Error reading request body")
		c.Abort()
		return
	}

	// Создание короткой ссылки
	shortURL, err := h.Service.CreateShortURL(string(body))
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating short URL")
		c.Abort()
		return
	}

	// Response: текст с полным URL
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, h.Configuration.ShortAddress+"/"+shortURL)
}

// Обработка POST запроса: сокращение URL (JSON)
func (h *Handler) SendJSONURL(c *gin.Context) {

	// Проверка Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		c.String(http.StatusBadRequest, "Invalid ContentType, application/json only")
		c.Abort()
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// Валидация реквеста
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// Создание короткой ссылки
	shortURL, err := h.Service.CreateShortURL(request.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating short URL"})
		c.Abort()
		return
	}

	// Response: JSON с полным URL
	var response model.Response
	response.Result = h.Configuration.ShortAddress + "/" + shortURL

	// Пишем ответ
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
}

// Обработка GET запроса: редирект по короткой ссылке
func (h *Handler) GetURL(c *gin.Context) {
	// Получаем параметр из URL: /:id
	shortURL := c.Param("id")

	// Ищем полную ссылку
	fullURL, err := h.Service.GetFullURL(shortURL)
	if err != nil {
		c.String(http.StatusBadRequest, "URL not found")
		c.Abort()
		return
	}

	// Редирект (307)
	c.Redirect(http.StatusTemporaryRedirect, fullURL)
}

// Ping PostgreSQL
func (h *Handler) Ping(c *gin.Context) {
	if err := h.Service.PingPostgreSQL(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
