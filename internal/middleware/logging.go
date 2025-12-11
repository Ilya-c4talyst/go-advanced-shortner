package middleware

import (
	"go.uber.org/zap"
)

// Sugar является глобальным экземпляром логгера SugaredLogger из библиотеки zap.
// Используется для структурированного логирования в приложении.
// Должен быть инициализирован вызовом InitLogger() перед использованием.
var Sugar zap.SugaredLogger

// InitLogger инициализирует глобальный логгер Sugar.
// Создает экземпляр development-логгера из библиотеки zap
// и настраивает его как SugaredLogger для удобного логирования.
//
// Паникует, если не удалось создать логгер.
//
// Пример использования:
//
//	middleware.InitLogger()
//	middleware.Sugar.Infoln("Application started")
func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	Sugar = *logger.Sugar()
}
