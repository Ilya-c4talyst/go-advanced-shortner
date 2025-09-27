package repository

// URLRepository интерфейс для работы с URL
type URLRepository interface {
	// GetFullValue получает оригинальный URL по короткому
	GetFullValue(shortURL string) (string, error)
	// GetShortValue получает короткий URL по оригинальному
	GetShortValue(shortURL string) (string, error)
	// SetValue сохраняет пару короткий URL - оригинальный URL с user_id
	SetValue(shortURL, originalURL, userID string) error
	// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL с user_id
	SetValuesBatch(pairs map[string]string, userID string) error
	// GetUserURLs получает все URL пользователя
	GetUserURLs(userID string) ([]map[string]string, error)
	// DeleteURLsBatch помечает URL как удаленные для указанного пользователя
	DeleteURLsBatch(shortURLs []string, userID string) error
	// IsDeleted проверяет, помечен ли URL как удаленный
	IsDeleted(shortURL string) (bool, error)
	// Close закрывает соединение с хранилищем
	Close() error
}
