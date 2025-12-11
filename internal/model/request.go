package model

// Request представляет запрос на создание короткого URL.
// Используется в эндпоинте POST /api/shorten.
type Request struct {
	URL string `json:"url" validate:"required,url"`
}

// Response представляет ответ с коротким URL.
// Используется в эндпоинте POST /api/shorten.
type Response struct {
	Result string `json:"result"`
}

// URLRecord представляет запись URL в JSON-хранилище.
// Используется для сериализации данных в файловом хранилище.
type URLRecord struct {
	ID          int    `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
}

// BatchRequest представляет один элемент пакетного запроса на создание коротких URL.
// Используется в эндпоинте POST /api/shorten/batch.
type BatchRequest struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	OriginalURL   string `json:"original_url" validate:"required,url"`
}

// BatchResponse представляет один элемент ответа на пакетный запрос.
// Используется в эндпоинте POST /api/shorten/batch.
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserURL представляет URL пользователя в ответе GET /api/user/urls.
// Содержит как короткий, так и оригинальный URL.
type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// DeleteURLsRequest представляет запрос на удаление URL пользователя.
// Используется в эндпоинте DELETE /api/user/urls.
// Представляет собой массив коротких URL для удаления.
type DeleteURLsRequest []string
