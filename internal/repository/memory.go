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

// withReadLock выполняет функцию под блокировкой чтения
func (r *MemoryRepository) withReadLock(fn func() error) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return fn()
}

// withWriteLock выполняет функцию под блокировкой записи
func (r *MemoryRepository) withWriteLock(fn func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return fn()
}

// checkURLExists проверяет существование и статус удаления URL (должен вызываться под блокировкой)
func (r *MemoryRepository) checkURLExists(shortURL string) (string, error) {
	// Проверяем, существует ли URL
	value, exists := r.data[shortURL]
	if !exists {
		return "", errors.New("not found key in database")
	}

	// Проверяем, не удален ли URL
	if deleted, exists := r.deletedMap[shortURL]; exists && deleted {
		return "", errors.New("not found key in database")
	}

	return value, nil
}

// GetValue получает оригинальный URL по короткому
func (r *MemoryRepository) GetFullValue(shortURL string) (string, error) {
	var result string
	var err error
	
	r.withReadLock(func() error {
		result, err = r.checkURLExists(shortURL)
		return err
	})
	
	return result, err
}

// GetShortValue получает короткий URL по оригинальному
func (r *MemoryRepository) GetShortValue(originalURL string) (string, error) {
	var result string
	var err error
	
	r.withReadLock(func() error {
		for short, long := range r.data {
			if long == originalURL {
				result = short
				return nil
			}
		}
		err = errors.New("not found key in database")
		return err
	})
	
	return result, err
}

// SetValue сохраняет пару короткий URL - оригинальный URL с user_id
func (r *MemoryRepository) SetValue(shortURL, originalURL, userID string) error {
	return r.withWriteLock(func() error {
		if _, ok := r.data[shortURL]; ok {
			return ErrRowExists
		}
		r.data[shortURL] = originalURL
		r.userMap[shortURL] = userID
		return nil
	})
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL с user_id
func (r *MemoryRepository) SetValuesBatch(pairs map[string]string, userID string) error {
	return r.withWriteLock(func() error {
		for key, value := range pairs {
			if _, ok := r.data[key]; ok {
				return ErrRowExists
			}
			r.data[key] = value
			r.userMap[key] = userID
		}
		return nil
	})
}

// Close закрывает соединение с хранилищем (для памяти это заглушка)
func (r *MemoryRepository) Close() error {
	return nil
}

// GetUserURLs получает все URL пользователя (не удаленные)
func (r *MemoryRepository) GetUserURLs(userID string) ([]map[string]string, error) {
	var urls []map[string]string
	
	r.withReadLock(func() error {
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
		return nil
	})
	
	return urls, nil
}

// DeleteURLsBatch помечает множественные URL как удаленные для указанного пользователя
func (r *MemoryRepository) DeleteURLsBatch(shortURLs []string, userID string) error {
	return r.withWriteLock(func() error {
		for _, shortURL := range shortURLs {
			// Проверяем, принадлежит ли URL пользователю
			if ownerID, exists := r.userMap[shortURL]; exists && ownerID == userID {
				r.deletedMap[shortURL] = true
			}
		}
		return nil
	})
}

// IsDeleted проверяет, помечен ли URL как удаленный
func (r *MemoryRepository) IsDeleted(shortURL string) (bool, error) {
	var result bool
	var err error
	
	r.withReadLock(func() error {
		// Проверяем, существует ли URL вообще
		if _, exists := r.data[shortURL]; !exists {
			err = errors.New("not found key in database")
			return err
		}

		// Проверяем статус удаления
		if deleted, exists := r.deletedMap[shortURL]; exists {
			result = deleted
		} else {
			result = false
		}
		return nil
	})
	
	return result, err
}
