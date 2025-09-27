package repository

import (
	"errors"
	"sync"
)

// MemoryRepository реализация репозитория для хранения в памяти
type MemoryRepository struct {
	data       map[string]string // shortURL -> originalURL
	userMap    map[string]string // shortURL -> userID
	deletedMap map[string]bool   // shortURL -> isDeleted
	mu         sync.RWMutex
}

// NewMemoryRepository создает новый репозиторий для работы с памятью
func NewMemoryRepository() URLRepository {
	return &MemoryRepository{
		data:       make(map[string]string),
		userMap:    make(map[string]string),
		deletedMap: make(map[string]bool),
	}
}

// GetValue получает оригинальный URL по короткому
func (r *MemoryRepository) GetFullValue(shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Проверяем, не удален ли URL
	if deleted, exists := r.deletedMap[shortURL]; exists && deleted {
		return "", errors.New("not found key in database")
	}

	if value, ok := r.data[shortURL]; ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// GetShortValue получает короткий URL по оригинальному
func (r *MemoryRepository) GetShortValue(originalURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for short, long := range r.data {
		if long == originalURL {
			return short, nil
		}
	}
	return "", errors.New("not found key in database")
}

// SetValue сохраняет пару короткий URL - оригинальный URL с user_id
func (r *MemoryRepository) SetValue(shortURL, originalURL, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[shortURL]; ok {
		return ErrRowExists
	}
	r.data[shortURL] = originalURL
	r.userMap[shortURL] = userID
	return nil
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL с user_id
func (r *MemoryRepository) SetValuesBatch(pairs map[string]string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, value := range pairs {
		if _, ok := r.data[key]; ok {
			return ErrRowExists
		}
		r.data[key] = value
		r.userMap[key] = userID
	}
	return nil
}

// Close закрывает соединение с хранилищем (для памяти это заглушка)
func (r *MemoryRepository) Close() error {
	return nil
}

// GetUserURLs получает все URL пользователя (не удаленные)
func (r *MemoryRepository) GetUserURLs(userID string) ([]map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var urls []map[string]string
	for shortURL, originalURL := range r.data {
		// Проверяем, что URL принадлежит пользователю и не удален
		if userID == r.userMap[shortURL] {
			if deleted, exists := r.deletedMap[shortURL]; !exists || !deleted {
				urls = append(urls, map[string]string{
					"short_url":    shortURL,
					"original_url": originalURL,
				})
			}
		}
	}
	return urls, nil
}

// DeleteURLsBatch помечает множественные URL как удаленные для указанного пользователя
func (r *MemoryRepository) DeleteURLsBatch(shortURLs []string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, shortURL := range shortURLs {
		// Проверяем, принадлежит ли URL пользователю
		if ownerID, exists := r.userMap[shortURL]; exists && ownerID == userID {
			r.deletedMap[shortURL] = true
		}
	}
	return nil
}

// IsDeleted проверяет, помечен ли URL как удаленный
func (r *MemoryRepository) IsDeleted(shortURL string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Проверяем, существует ли URL вообще
	if _, exists := r.data[shortURL]; !exists {
		return false, errors.New("not found key in database")
	}

	// Проверяем статус удаления
	if deleted, exists := r.deletedMap[shortURL]; exists {
		return deleted, nil
	}
	return false, nil
}
