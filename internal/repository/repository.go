package repository

// URLRepository интерфейс для работы с URL
type URLRepository interface {
	// GetValue получает оригинальный URL по короткому
	GetValue(shortURL string) (string, error)
	// SetValue сохраняет пару короткий URL - оригинальный URL
	SetValue(shortURL, originalURL string) error
	// Close закрывает соединение с хранилищем
	Close() error
}
