package repository

import (
	"errors"
	"sync"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/persistence"
)

// FileRepository реализация репозитория для хранения в файле
type FileRepository struct {
	data        map[string]string
	mu          sync.RWMutex
	filePath    string
	persistence persistence.JSONPersistence
}

// NewFileRepository создает новый репозиторий для работы с файлом
func NewFileRepository(filePath string) URLRepository {
	repo := &FileRepository{
		data:        make(map[string]string),
		filePath:    filePath,
		persistence: persistence.NewFileJSONPersistence(),
	}
	
	// Загружаем данные из файла при инициализации
	data, _, err := repo.persistence.Load(filePath)
	if err == nil {
		repo.data = data
	}
	
	return repo
}

// GetValue получает оригинальный URL по короткому
func (r *FileRepository) GetValue(shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if value, ok := r.data[shortURL]; ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// SetValue сохраняет пару короткий URL - оригинальный URL
func (r *FileRepository) SetValue(shortURL, originalURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.data[shortURL] = originalURL
	
	// Сохраняем в файл
	return r.persistence.Save(r.filePath, r.data)
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL
func (r *FileRepository) SetValuesBatch(pairs map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	for shortURL, originalURL := range pairs {
		r.data[shortURL] = originalURL
	}
	
	// Сохраняем в файл
	return r.persistence.Save(r.filePath, r.data)
}

// Close закрывает соединение с хранилищем
func (r *FileRepository) Close() error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Финальное сохранение в файл
	return r.persistence.Save(r.filePath, r.data)
}
