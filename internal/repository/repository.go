package repository

// URLReader определяет интерфейс для чтения URL из хранилища.
// Интерфейс предоставляет методы для получения оригинальных URL,
// коротких URL и списка URL пользователя.
type URLReader interface {
	// GetFullValue получает оригинальный URL по короткому идентификатору.
	// Параметры:
	//   - shortURL: короткий URL (без базового адреса)
	// Возвращает:
	//   - string: оригинальный URL
	//   - error: ошибка, если URL не найден
	GetFullValue(shortURL string) (string, error)

	// GetShortValue получает короткий URL по оригинальному URL.
	// Используется для проверки существования URL.
	// Параметры:
	//   - shortURL: оригинальный URL
	// Возвращает:
	//   - string: короткий URL
	//   - error: ошибка, если URL не найден
	GetShortValue(shortURL string) (string, error)

	// GetUserURLs получает все URL, созданные указанным пользователем.
	// Параметры:
	//   - userID: идентификатор пользователя
	// Возвращает:
	//   - []map[string]string: массив карт с ключами "short_url" и "original_url"
	//   - error: ошибка при получении URL
	GetUserURLs(userID string) ([]map[string]string, error)

	// IsDeleted проверяет, помечен ли URL как удаленный (soft delete).
	// Параметры:
	//   - shortURL: короткий URL для проверки
	// Возвращает:
	//   - bool: true если URL удален, false в противном случае
	//   - error: ошибка при проверке статуса
	IsDeleted(shortURL string) (bool, error)
}

// URLWriter определяет интерфейс для записи и модификации URL в хранилище.
// Интерфейс предоставляет методы для создания и удаления URL.
type URLWriter interface {
	// SetValue сохраняет связь между коротким и оригинальным URL.
	// Параметры:
	//   - shortURL: короткий URL (без базового адреса)
	//   - originalURL: оригинальный URL
	//   - userID: идентификатор пользователя, создающего URL
	// Возвращает:
	//   - error: ошибка при сохранении или ErrRowExists если URL уже существует
	SetValue(shortURL, originalURL, userID string) error

	// SetValuesBatch сохраняет пакет связей между короткими и оригинальными URL.
	// Более эффективен для больших объемов данных.
	// Параметры:
	//   - pairs: карта соответствия коротких URL к оригинальным
	//   - userID: идентификатор пользователя, создающего URL
	// Возвращает:
	//   - error: ошибка при сохранении
	SetValuesBatch(pairs map[string]string, userID string) error

	// DeleteURLsBatch помечает URL как удаленные (soft delete) для указанного пользователя.
	// Параметры:
	//   - shortURLs: массив коротких URL для удаления
	//   - userID: идентификатор пользователя, владеющего URL
	// Возвращает:
	//   - error: ошибка при удалении URL
	DeleteURLsBatch(shortURLs []string, userID string) error
}

// URLRepository объединяет интерфейсы URLReader и URLWriter,
// предоставляя полный набор операций для работы с URL-хранилищем.
type URLRepository interface {
	URLReader
	URLWriter

	// Close закрывает соединение с хранилищем и освобождает ресурсы.
	// Должен быть вызван при завершении работы приложения.
	// Возвращает:
	//   - error: ошибка при закрытии соединения
	Close() error
}
