package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UserIDKey является ключом для хранения ID пользователя в контексте Gin.
// Используется для передачи идентификатора пользователя между middleware и хендлерами.
const UserIDKey = "userID"

// LoggingMiddleware создает middleware для логирования HTTP-запросов.
// Middleware записывает информацию о каждом запросе: URI, метод, статус, длительность и IP клиента.
//
// Возвращает:
//   - gin.HandlerFunc: функция middleware для логирования
//
// Пример использования:
//
//	router.Use(middleware.LoggingMiddleware())
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Обслуживаем запрос
		c.Next()

		duration := time.Since(start)

		Sugar.Infoln(
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", duration,
			"client_ip", c.ClientIP(),
		)
	}
}

// GzipMiddleware создает middleware для обработки gzip сжатия запросов и ответов.
// Middleware автоматически распаковывает входящие gzip-сжатые данные
// и сжимает исходящие ответы, если клиент поддерживает gzip.
//
// Возвращает:
//   - gin.HandlerFunc: функция middleware для обработки gzip
//
// Пример использования:
//
//	router.Use(middleware.GzipMiddleware())
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Обрабатываем входящие сжатые данные
		if c.GetHeader("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid gzip data"})
				return
			}
			defer reader.Close()

			// Заменяем тело запроса на распакованное
			body, err := io.ReadAll(reader)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error reading gzip data"})
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Request.ContentLength = int64(len(body))
			c.Request.Header.Del("Content-Encoding")

			// Восстанавливаем правильный Content-Type для разных случаев
			contentType := c.GetHeader("Content-Type")
			if contentType == "application/x-gzip" {
				// Определяем тип содержимого по тому, что делаем
				if strings.Contains(c.Request.URL.Path, "/api/") {
					c.Request.Header.Set("Content-Type", "application/json")
				} else {
					c.Request.Header.Set("Content-Type", "text/plain")
				}
			}
		}

		// Проверяем, поддерживает ли клиент gzip для ответа
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Устанавливаем заголовки для сжатого ответа
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Создаем gzip writer для ответа
		gz := gzip.NewWriter(c.Writer)
		defer func() {
			gz.Close()
		}()

		// Оборачиваем ResponseWriter
		c.Writer = &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}

		c.Next()
	}
}

// AuthServiceInterface определяет интерфейс для сервиса аутентификации.
// Реализация должна предоставлять методы для работы с пользовательскими cookie
// и аутентификацией пользователей.
type AuthServiceInterface interface {
	// GetOrCreateUserID извлекает ID пользователя из запроса или создает нового пользователя
	GetOrCreateUserID(r *http.Request) (string, *http.Cookie)

	// ValidateCookie проверяет подлинность cookie и возвращает ID пользователя
	ValidateCookie(cookieValue string) (string, bool)
}

// AuthMiddleware создает middleware для аутентификации пользователей.
// Middleware получает или создает ID пользователя из cookie и сохраняет его в контексте запроса.
// Если пользователь не аутентифицирован, создается новый пользователь и устанавливается cookie.
//
// Параметры:
//   - authService: сервис аутентификации, реализующий AuthServiceInterface
//
// Возвращает:
//   - gin.HandlerFunc: функция middleware для обработки аутентификации
//
// Пример использования:
//
//	router.Use(middleware.AuthMiddleware(authService))
func AuthMiddleware(authService AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, cookie := authService.GetOrCreateUserID(c.Request)

		// Сохраняем userID в контексте для использования в хендлерах
		c.Set(UserIDKey, userID)

		// Если была создана новая кука, устанавливаем её
		if cookie != nil {
			http.SetCookie(c.Writer, cookie)
		}

		c.Next()
	}
}

// RequireAuthMiddleware создает middleware, который требует валидную существующую аутентификацию.
// В отличие от AuthMiddleware, этот middleware не создает нового пользователя,
// а возвращает ошибку 401 Unauthorized, если cookie отсутствует или невалидна.
//
// Параметры:
//   - authService: сервис аутентификации, реализующий AuthServiceInterface
//
// Возвращает:
//   - gin.HandlerFunc: функция middleware для строгой проверки аутентификации
//
// Пример использования:
//
//	router.GET("/api/user/urls", middleware.RequireAuthMiddleware(authService), handler.GetUserURLs)
func RequireAuthMiddleware(authService AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем наличие куки в запросе
		cookie, err := c.Request.Cookie("user_id")
		if err != nil {
			// Кука отсутствует
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Проверяем валидность куки
		userID, valid := authService.ValidateCookie(cookie.Value)
		if !valid || userID == "" {
			// Кука невалидна
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Сохраняем userID в контексте для использования в хендлерах
		c.Set(UserIDKey, userID)

		c.Next()
	}
}
