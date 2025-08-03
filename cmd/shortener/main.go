package main

import (
	"log"
	"net/http"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/handler"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/storage"
)

func main() {
	// Инициализация роутера
	mux := http.NewServeMux()

	// Создание БД(пока фейк) и репозитория
	db := storage.CreateDB()
	repo := repository.NewShortenerRepository(db)

	// Создание сервиса
	shortService := service.NewURLShortnerService(repo)

	// Создание обработчика
	handler.NewHandler(mux, shortService)

	// Запуск сервера
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
