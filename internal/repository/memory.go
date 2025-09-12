package repository

import (
	"maps"
	"errors"
	"sync"
)

// MemoryRepository реализация репозитория для хранения в памяти
type MemoryRepository struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewMemoryRepository создает новый репозиторий для работы с памятью
func NewMemoryRepository() URLRepository {
	return &MemoryRepository{
		data: make(map[string]string),
	}
}

// GetValue получает оригинальный URL по короткому
func (r *MemoryRepository) GetValue(shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if value, ok := r.data[shortURL]; ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// SetValue сохраняет пару короткий URL - оригинальный URL
func (r *MemoryRepository) SetValue(shortURL, originalURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.data[shortURL] = originalURL
	return nil
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL
func (r *MemoryRepository) SetValuesBatch(pairs map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	maps.Copy(r.data, pairs)
	return nil
}

// Close закрывает соединение с хранилищем (для памяти это заглушка)
func (r *MemoryRepository) Close() error {
	return nil
}
