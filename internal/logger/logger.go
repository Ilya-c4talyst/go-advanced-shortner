package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Singleton для логгера
var sugar zap.SugaredLogger

// Инициализатор для логгера
func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()
}

// Middleware для логирования запросов
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Обслуживаем запрос
		c.Next()

		duration := time.Since(start)

		sugar.Infoln(
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", duration,
			"client_ip", c.ClientIP(),
		)
	}
}
