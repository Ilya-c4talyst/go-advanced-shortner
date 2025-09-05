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
	DB            *repository.ShortenerRepository
	Configuration *config.ConfigStruct
}

// Конструктор для сервиса
func NewURLShortnerService(repo *repository.ShortenerRepository, configuration *config.ConfigStruct) *URLShortnerService {

	return &URLShortnerService{
		DB:            repo,
		Configuration: configuration,
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

// Ping DB
func (u *URLShortnerService) PingPostgreSQL() error {
	db, err := sql.Open("pgx", u.Configuration.AddressDB)
	if err != nil {
		return err
	}
	err = db.Ping()
	return err
}
