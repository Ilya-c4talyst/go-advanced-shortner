package main

import (
	"log"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/handler"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/middleware"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/server"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Чтение флагов
	configuration := config.GenerateConfig()

	// Инициализация логгера
	middleware.InitLogger()

	// Инициализация роутера
	ginEngine := gin.Default()

	// Создание БД и загрузка данных из файла
	db := storage.CreateDB()
	if err := db.LoadFromFile(configuration.FilePath); err != nil {
		log.Fatalf("Ошибка загрузки данных из файла: %v", err)
	}

	// Создание репозитория
	repo := repository.NewShortenerRepository(db)

	// Создание сервиса
	shortService := service.NewURLShortnerService(repo)
	
	// Создание обработчика
	handler.NewHandler(ginEngine, shortService, configuration)

	// Создание и запуск сервера с graceful shutdown
	srv := server.NewServer(configuration.Port, ginEngine, db, configuration.FilePath)
	if err := srv.Start(); err != nil {
		log.Fatalf("Ошибка работы сервера: %v", err)
	}

	log.Println("Сервер завершён")
}
