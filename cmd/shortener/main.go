package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Создание HTTP сервера
	server := &http.Server{
		Addr:    configuration.Port,
		Handler: ginEngine,
	}

	// Запуск сервера в горутине
	go func() {
		log.Printf("Сервер запущен на %s", configuration.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при завершении сервера: %v", err)
	}

	// Сохранение данных перед завершением работы
	log.Println("Сохранение данных...")
	if err := db.SaveToFile(configuration.FilePath); err != nil {
		log.Printf("Ошибка сохранения данных: %v", err)
	} else {
		log.Println("Данные успешно сохранены")
	}

	log.Println("Сервер завершён")
}
