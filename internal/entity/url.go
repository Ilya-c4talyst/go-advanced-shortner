package entity

// URL представляет доменную модель URL
type URL struct {
	ID          int
	ShortURL    string
	OriginalURL string
	UserID      string
	IsDeleted   bool
}

// NewURL создаёт новую доменную сущность URL
func NewURL(id int, shortURL, originalURL, userID string) *URL {
	return &URL{
		ID:          id,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
		IsDeleted:   false,
	}
}

// GetShortURL возвращает сокращённый URL
func (u *URL) GetShortURL() string {
	return u.ShortURL
}

// GetOriginalURL возвращает оригинальный URL
func (u *URL) GetOriginalURL() string {
	return u.OriginalURL
}

// GetID возвращает идентификатор
func (u *URL) GetID() int {
	return u.ID
}

// IsURLDeleted возвращает статус удаления URL
func (u *URL) IsURLDeleted() bool {
	return u.IsDeleted
}

// MarkAsDeleted помечает URL как удаленный
func (u *URL) MarkAsDeleted() {
	u.IsDeleted = true
}
