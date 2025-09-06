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
		if _, err := u.Repository.GetValue(shortURL); err == nil {
			continue
		}
		break
	}

	// Сохраняем в репозитории
	if err := u.Repository.SetValue(shortURL, url); err != nil {
		return "", err
	}
	
	return shortURL, nil
}

// Получение полного URL
func (u *URLShortnerService) GetFullURL(shortURL string) (string, error) {
	// Ищем полный URL в репозитории, или выдаем ошибку
	if url, err := u.Repository.GetValue(shortURL); err == nil {
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
