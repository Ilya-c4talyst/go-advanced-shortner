package main

import (
	"log"
	"net/http"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/handler"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/middleware"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
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

	// Создание БД(пока фейк) и репозитория
	db := storage.CreateDB()
	repo := repository.NewShortenerRepository(db)

	// Создание сервиса
	shortService := service.NewURLShortnerService(repo)

	// Создание обработчика
	handler.NewHandler(ginEngine, shortService, configuration)

	// Загрузка данных из файла
	err := db.LoadFromFile(configuration.FilePath)
	if err != nil {
		log.Fatalf("Ошибка загрузки данных из файла: %v", err)
	}

	defer func() {
		log.Println("Сохранение данных перед завершением работы...")
		if err := db.SaveToFile(configuration.FilePath); err != nil {
			log.Printf("Ошибка сохранения данных: %v", err)
		}
	}()

	// Запуск сервера
	err = http.ListenAndServe(configuration.Port, ginEngine)
	if err != nil {
		log.Fatal(err)
	}
}
