package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
)

// Структура хендера
type Handler struct {
	Service *service.URLShortnerService
}

// Конструктор для хендлера
func NewHandler(mux *http.ServeMux, service *service.URLShortnerService) {
	handler := &Handler{
		service,
	}
	mux.HandleFunc("POST /", handler.SendURL())
	mux.HandleFunc("GET /{id}", handler.GetURL())
}

// Обработка POST запроса
func (h *Handler) SendURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Проверка на тип контента
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(strings.ToLower(contentType), "text/plain") {
			http.Error(w, "Invalid ContentType, text/plain only", http.StatusBadRequest)
			return
		}

		// Чтение тела запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Получаем сокращенную ссылку из сервиса
		shortUrl := h.Service.CreateShortURL(string(body))

		// Пишем ответ в респонс
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + shortUrl))
	}
}

func (h *Handler) GetURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Получаем ссылку из пути
		shortUrl := r.PathValue("id")

		// Ищем ссылку в БД
		fullUrl, err := h.Service.GetFullURL(shortUrl)

		// Обрабатываем ошибку, если не нашли URL
		if err != nil {
			http.Error(w, "Error not found", http.StatusBadRequest)
			return
		}

		// Пишем ответ в респонс
		w.Header().Set("Location", fullUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		// Почему не проходит этот вариант...
		// http.Redirect(w, r, fullUrl, http.StatusTemporaryRedirect)
	}
}
