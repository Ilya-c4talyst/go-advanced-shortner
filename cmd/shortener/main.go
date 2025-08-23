package main

import (
	"log"
	"net/http"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/handler"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/logger"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Чтение флагов
	configuration := config.GenerateConfig()

	// Инициализация логгера
	logger.InitLogger()

	// Инициализация роутера
	ginEngine := gin.Default()

	// Создание БД(пока фейк) и репозитория
	db := storage.CreateDB()
	repo := repository.NewShortenerRepository(db)

	// Создание сервиса
	shortService := service.NewURLShortnerService(repo)

	// Создание обработчика
	handler.NewHandler(ginEngine, shortService, configuration)

	// Запуск сервера
	err := http.ListenAndServe(configuration.Port, ginEngine)
	if err != nil {
		log.Fatal(err)
	}
}
