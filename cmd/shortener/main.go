package main

import (
	"log"
	"net/http"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/handler"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
)

func main() {
	// Инициализация роутера
	mux := http.NewServeMux()

	// Создание сервиса
	shortService := service.NewUrlShortnerService()

	// Создание обработчика
	handler.NewHandler(mux, shortService)

	// Запуск сервера
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
