package main

import (
	"log"
	"net/http"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/handler"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация роутера
	ginEngine := gin.Default()

	// Создание БД(пока фейк) и репозитория
	db := storage.CreateDB()
	repo := repository.NewShortenerRepository(db)

	// Создание сервиса
	shortService := service.NewURLShortnerService(repo)

	// Создание обработчика
	handler.NewHandler(ginEngine, shortService)

	// Запуск сервера
	err := http.ListenAndServe(config.Configuration.Port, ginEngine)
	if err != nil {
		log.Fatal(err)
	}
}
