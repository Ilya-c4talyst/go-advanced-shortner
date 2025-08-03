package repository

import (
	"testing"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestShortenerRepository(t *testing.T) {
	// Инициализация тестовой БД и репозитория
	db := storage.CreateDB()
	repo := NewShortenerRepository(db)

	t.Run("Set and Get value successfully", func(t *testing.T) {
		key := "testKey"
		value := "https://example.com"

		// Записываем значение
		repo.SetValue(key, value)

		// Получаем значение и проверяем
		result, err := repo.GetValue(key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Get non-existent key returns error", func(t *testing.T) {
		nonExistentKey := "nonExistentKey"

		// Пытаемся получить несуществующий ключ
		result, err := repo.GetValue(nonExistentKey)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Equal(t, "not found key in database", err.Error())
	})

	t.Run("Overwrite existing key", func(t *testing.T) {
		key := "existingKey"
		firstValue := "https://first.com"
		secondValue := "https://second.com"

		// Первая запись
		repo.SetValue(key, firstValue)
		firstResult, err := repo.GetValue(key)
		assert.NoError(t, err)
		assert.Equal(t, firstValue, firstResult)

		// Перезаписываем
		repo.SetValue(key, secondValue)
		secondResult, err := repo.GetValue(key)
		assert.NoError(t, err)
		assert.Equal(t, secondValue, secondResult)
	})
}
