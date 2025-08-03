package service

import (
	"errors"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/pkg/utils"
)

// Структура для сервиса сокращения ссылок
type URLShortnerService struct {
	DB *repository.ShortenerRepository
}

// Конструктор для сервиса
func NewURLShortnerService(repo *repository.ShortenerRepository) *URLShortnerService {

	return &URLShortnerService{
		DB: repo,
	}
}

// Создание сокращенного URL
func (u *URLShortnerService) CreateShortURL(url string) string {

	// Инициализация результата
	var shortURL string

	// Генерируем сокращенную уникальную сслыку
	for {
		shortURL = utils.GenerateShortKey()
		if _, err := u.DB.GetValue(shortURL); err == nil {
			continue
		}
		break
	}

	// Сохраняем в БД
	u.DB.SetValue(shortURL, url)
	return shortURL
}

// Получение полного URL
func (u *URLShortnerService) GetFullURL(shortURL string) (string, error) {

	// Ищем полный URL в БД, или выдаем ошибку
	if url, err := u.DB.GetValue(shortURL); err == nil {
		return url, nil
	} else {
		return "", errors.New("not found")
	}
}
