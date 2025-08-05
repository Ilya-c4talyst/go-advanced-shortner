package handler

import (
	"net/http"
	"strings"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/gin-gonic/gin"
)

// Handler — структура хендлера
type Handler struct {
	Service *service.URLShortnerService
}

// Конструктор для хендлера
func NewHandler(ginEngine *gin.Engine, service *service.URLShortnerService) {
	handler := &Handler{Service: service}

	// Регистрируем маршруты
	ginEngine.POST("/", handler.SendURL)
	ginEngine.GET("/:id", handler.GetURL)
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
	shortURL := h.Service.CreateShortURL(string(body))

	// Response: текст с полным URL
	c.Header("Content-Type", "text/plain")
	c.Status(http.StatusCreated)
	c.String(http.StatusCreated, "http://localhost:8080/"+shortURL)
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
