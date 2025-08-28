package repository

import (
	"errors"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/storage"
)

// Структура для репозитория
type ShortenerRepository struct {
	Database *storage.DB
}

// Конструктор для репозитория
func NewShortenerRepository(database *storage.DB) *ShortenerRepository {
	return &ShortenerRepository{
		Database: database,
	}
}

// Получение данных из БД
func (r *ShortenerRepository) GetValue(key string) (string, error) {
	if value, ok := r.Database.Get(key); ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// Запись данных в БД
func (r *ShortenerRepository) SetValue(key string, value string) {
	r.Database.Set(key, value)
}
