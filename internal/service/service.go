package service

import (
	"database/sql"
	"errors"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/pkg/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Структура для сервиса сокращения ссылок
type URLShortnerService struct {
	Repository    repository.URLRepository
	Configuration *config.ConfigStruct
}

// Конструктор для сервиса
func NewURLShortnerService(repo repository.URLRepository, configuration *config.ConfigStruct) *URLShortnerService {
	return &URLShortnerService{
		Repository:    repo,
		Configuration: configuration,
	}
}

// Создание сокращенного URL
func (u *URLShortnerService) CreateShortURL(url string) (string, error) {
	// Инициализация результата
	var shortURL string

	// Генерируем сокращенную уникальную ссылку
	for {
		shortURL = utils.GenerateShortKey()
		if _, err := u.Repository.GetFullValue(shortURL); err == nil {
			continue
		}
		break
	}

	// Сохраняем в репозитории
	if err := u.Repository.SetValue(shortURL, url); err != nil {
		if errors.Is(err, repository.ErrRowExists) {
			// Если ссылка уже существует
			if shortURL, err = u.Repository.GetShortValue(url); err == nil {
				return shortURL, repository.ErrRowExists
			}
		}
		return "", err
	}

	return shortURL, nil
}

// CreateShortURLsBatch создает сокращенные URL для пакета URL
func (u *URLShortnerService) CreateShortURLsBatch(urls []string) (map[string]string, error) {
	if len(urls) == 0 {
		return make(map[string]string), nil
	}

	result := make(map[string]string)
	pairs := make(map[string]string)

	// Генерируем короткие URL для каждого исходного URL
	for _, originalURL := range urls {
		var shortURL string

		// Генерируем уникальную короткую ссылку
		for {
			shortURL = utils.GenerateShortKey()
			if _, err := u.Repository.GetFullValue(shortURL); err == nil {
				continue
			}
			// Проверяем также, что этот ключ не используется в текущем пакете
			if _, exists := pairs[shortURL]; exists {
				continue
			}
			break
		}

		pairs[shortURL] = originalURL
		result[originalURL] = shortURL
	}

	// Сохраняем пакет в репозитории
	if err := u.Repository.SetValuesBatch(pairs); err != nil {
		return nil, err
	}

	return result, nil
}

// Получение полного URL
func (u *URLShortnerService) GetFullURL(shortURL string) (string, error) {
	// Ищем полный URL в репозитории, или выдаем ошибку
	if url, err := u.Repository.GetFullValue(shortURL); err == nil {
		return url, nil
	} else {
		return "", errors.New("not found")
	}
}

// Ping DB
func (u *URLShortnerService) PingPostgreSQL() error {
	db, err := sql.Open("pgx", u.Configuration.AddressDB)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

// Close закрывает соединение с репозиторием
func (u *URLShortnerService) Close() error {
	return u.Repository.Close()
}
