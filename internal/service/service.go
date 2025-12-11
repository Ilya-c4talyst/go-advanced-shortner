package service

import (
	"database/sql"
	"errors"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/pkg/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// URLShortnerService предоставляет бизнес-логику для сервиса сокращения URL.
// Сервис управляет созданием, получением и удалением коротких URL,
// взаимодействуя с репозиторием для хранения данных.
type URLShortnerService struct {
	Repository    repository.URLRepository
	Configuration *config.ConfigStruct
}

// NewURLShortnerService создает новый экземпляр сервиса сокращения URL.
//
// Параметры:
//   - repo: репозиторий для работы с хранилищем URL
//   - configuration: конфигурация приложения
//
// Возвращает:
//   - *URLShortnerService: инициализированный экземпляр сервиса
func NewURLShortnerService(repo repository.URLRepository, configuration *config.ConfigStruct) *URLShortnerService {
	return &URLShortnerService{
		Repository:    repo,
		Configuration: configuration,
	}
}

// CreateShortURL создает короткий URL для указанного оригинального URL.
// Метод генерирует уникальный короткий ключ и сохраняет связь с оригинальным URL.
//
// Параметры:
//   - url: оригинальный URL для сокращения
//   - userID: идентификатор пользователя, создающего URL
//
// Возвращает:
//   - string: сгенерированный короткий URL (без базового адреса)
//   - error: ошибка при создании URL или repository.ErrRowExists если URL уже существует
//
// Пример:
//
//	shortURL, err := service.CreateShortURL("https://example.com/long-url", "user123")
//	if err != nil {
//	    if errors.Is(err, repository.ErrRowExists) {
//	        // URL уже существует, shortURL содержит существующий короткий URL
//	    }
//	    return err
//	}
//	// shortURL содержит новый короткий URL, например "abc123"
func (u *URLShortnerService) CreateShortURL(url, userID string) (string, error) {
	// Инициализация результата
	var shortURL string

	// Генерируем сокращенную уникальную ссылку
	for {
		shortURL = utils.GenerateShortKey()
		if _, err := u.Repository.GetFullValue(shortURL); err == nil {
			continue
		}
		break
	}

	// Сохраняем в репозитории
	if err := u.Repository.SetValue(shortURL, url, userID); err != nil {
		if errors.Is(err, repository.ErrRowExists) {
			// Если ссылка уже существует
			if shortURL, err = u.Repository.GetShortValue(url); err == nil {
				return shortURL, repository.ErrRowExists
			}
		}
		return "", err
	}

	return shortURL, nil
}

// CreateShortURLsBatch создает короткие URL для пакета оригинальных URL.
// Метод обрабатывает несколько URL одновременно, что более эффективно для больших объемов данных.
//
// Параметры:
//   - urls: массив оригинальных URL для сокращения
//   - userID: идентификатор пользователя, создающего URL
//
// Возвращает:
//   - map[string]string: карта соответствия оригинальных URL к коротким URL
//   - error: ошибка при создании URL
//
// Пример:
//
//	urls := []string{"https://example.com/url1", "https://example.com/url2"}
//	result, err := service.CreateShortURLsBatch(urls, "user123")
//	if err != nil {
//	    return err
//	}
//	// result["https://example.com/url1"] == "abc123"
//	// result["https://example.com/url2"] == "def456"
func (u *URLShortnerService) CreateShortURLsBatch(urls []string, userID string) (map[string]string, error) {
	if len(urls) == 0 {
		return make(map[string]string), nil
	}

	result := make(map[string]string)
	pairs := make(map[string]string)

	// Генерируем короткие URL для каждого исходного URL
	for _, originalURL := range urls {
		var shortURL string

		// Генерируем уникальную короткую ссылку
		for {
			shortURL = utils.GenerateShortKey()
			if _, err := u.Repository.GetFullValue(shortURL); err == nil {
				continue
			}
			// Проверяем также, что этот ключ не используется в текущем пакете
			if _, exists := pairs[shortURL]; exists {
				continue
			}
			break
		}

		pairs[shortURL] = originalURL
		result[originalURL] = shortURL
	}

	// Сохраняем пакет в репозитории
	if err := u.Repository.SetValuesBatch(pairs, userID); err != nil {
		return nil, err
	}

	return result, nil
}

// GetFullURL возвращает оригинальный URL по короткому URL.
// Метод проверяет, не был ли URL удален, прежде чем вернуть его.
//
// Параметры:
//   - shortURL: короткий URL (без базового адреса)
//
// Возвращает:
//   - string: оригинальный URL
//   - error: ошибка, если URL не найден или был удален
//
// Пример:
//
//	fullURL, err := service.GetFullURL("abc123")
//	if err != nil {
//	    // URL не найден или удален
//	    return err
//	}
//	// fullURL == "https://example.com/long-url"
func (u *URLShortnerService) GetFullURL(shortURL string) (string, error) {
	// Проверяем, не удален ли URL
	if deleted, err := u.Repository.IsDeleted(shortURL); err == nil && deleted {
		return "", errors.New("url is deleted")
	}

	// Ищем полный URL в репозитории, или выдаем ошибку
	if url, err := u.Repository.GetFullValue(shortURL); err == nil {
		return url, nil
	} else {
		return "", errors.New("not found")
	}
}

// GetUserURLs возвращает все URL, созданные указанным пользователем.
//
// Параметры:
//   - userID: идентификатор пользователя
//
// Возвращает:
//   - []map[string]string: массив карт с ключами "short_url" и "original_url"
//   - error: ошибка при получении URL
func (u *URLShortnerService) GetUserURLs(userID string) ([]map[string]string, error) {
	return u.Repository.GetUserURLs(userID)
}

// PingPostgreSQL проверяет доступность базы данных PostgreSQL.
// Метод используется для healthcheck приложения.
//
// Возвращает:
//   - error: ошибка, если база данных недоступна, или nil если соединение успешно
func (u *URLShortnerService) PingPostgreSQL() error {
	db, err := sql.Open("pgx", u.Configuration.AddressDB)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

// DeleteURLsBatch помечает указанные URL как удаленные для пользователя.
// Удаление происходит мягко (soft delete) - URL помечаются флагом is_deleted.
//
// Параметры:
//   - shortURLs: массив коротких URL для удаления
//   - userID: идентификатор пользователя, владеющего URL
//
// Возвращает:
//   - error: ошибка при удалении URL
func (u *URLShortnerService) DeleteURLsBatch(shortURLs []string, userID string) error {
	return u.Repository.DeleteURLsBatch(shortURLs, userID)
}

// IsDeleted проверяет, помечен ли URL как удаленный.
//
// Параметры:
//   - shortURL: короткий URL для проверки
//
// Возвращает:
//   - bool: true если URL удален, false в противном случае
//   - error: ошибка при проверке статуса
func (u *URLShortnerService) IsDeleted(shortURL string) (bool, error) {
	return u.Repository.IsDeleted(shortURL)
}

// Close закрывает соединение с репозиторием и освобождает ресурсы.
//
// Возвращает:
//   - error: ошибка при закрытии соединения
func (u *URLShortnerService) Close() error {
	return u.Repository.Close()
}
