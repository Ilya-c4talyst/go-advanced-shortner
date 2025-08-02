package service

import (
	"errors"

	"github.com/Ilya-c4talyst/go-advanced-shortner/pkg/utils"
)

// Структура для сервиса сокращения ссылок
type URLShortnerService struct {
	DB map[string]string
}

// Конструктор для сервиса
func NewURLShortnerService() *URLShortnerService {

	// Создание фековой БД
	db := make(map[string]string)

	return &URLShortnerService{
		DB: db,
	}
}

// Создание сокращенного URL
func (u *URLShortnerService) CreateShortURL(url string) string {

	// Инициализация результата
	var shortUrl string

	// Генерируем сокращенную уникальную сслыку
	for {
		shortUrl = utils.GenerateShortKey()
		if _, ok := u.DB[shortUrl]; ok {
			continue
		}
		break
	}

	// Сохраняем в БД
	u.DB[shortUrl] = url
	return shortUrl
}

// Получение полного URL
func (u *URLShortnerService) GetFullURL(shortUrl string) (string, error) {

	// Ищем полный URL в БД, или выдаем ошибку
	if url, ok := u.DB[shortUrl]; ok {
		return url, nil
	} else {
		return "", errors.New("not found")
	}
}
