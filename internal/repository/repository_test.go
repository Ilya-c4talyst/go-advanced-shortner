package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryRepository(t *testing.T) {
	// Инициализация репозитория в памяти
	repo := NewMemoryRepository()
	defer repo.Close()

	t.Run("Set and Get value successfully", func(t *testing.T) {
		key := "testKey"
		value := "https://example.com"
		userID := "user123"

		// Записываем значение
		err := repo.SetValue(key, value, userID)
		assert.NoError(t, err)

		// Получаем значение и проверяем
		result, err := repo.GetFullValue(key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Get non-existent key returns error", func(t *testing.T) {
		nonExistentKey := "nonExistentKey"

		// Пытаемся получить несуществующий ключ
		result, err := repo.GetFullValue(nonExistentKey)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Equal(t, "not found key in database", err.Error())
	})

	t.Run("Overwrite existing key", func(t *testing.T) {
		key := "existingKey"
		firstValue := "https://first.com"
		secondValue := "https://second.com"
		userID := "user456"

		// Первая запись
		err := repo.SetValue(key, firstValue, userID)
		assert.NoError(t, err)
		firstResult, err := repo.GetFullValue(key)
		assert.NoError(t, err)
		assert.Equal(t, firstValue, firstResult)

		// Перезаписываем
		err = repo.SetValue(key, secondValue, userID)
		assert.Error(t, err)
		_, err = repo.GetFullValue(key)
		assert.NoError(t, err)
	})

	t.Run("Get user URLs", func(t *testing.T) {
		userID := "user789"
		
		// Создаем несколько URL для пользователя
		_ = repo.SetValue("key1", "https://example1.com", userID)
		_ = repo.SetValue("key2", "https://example2.com", userID)
		_ = repo.SetValue("key3", "https://example3.com", "anotherUser")

		// Получаем URL пользователя
		userURLs, err := repo.GetUserURLs(userID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(userURLs))
	})
}
