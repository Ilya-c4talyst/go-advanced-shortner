package audit

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// AuditEvent представляет событие аудита в системе.
// Содержит информацию о действиях пользователей для целей мониторинга и анализа.
type AuditEvent struct {
	// Timestamp содержит Unix timestamp события в секундах
	Timestamp int64 `json:"ts"`

	// Action определяет тип действия (shorten или follow)
	Action string `json:"action"`

	// UserID содержит идентификатор пользователя, выполнившего действие
	UserID string `json:"user_id"`

	// URL содержит оригинальный (не сокращенный) URL
	URL string `json:"url"`
}

// Константы для типов действий в аудите
const (
	// ActionShorten представляет действие создания сокращенного URL
	ActionShorten = "shorten"

	// ActionFollow представляет действие перехода по короткой ссылке
	ActionFollow = "follow"
)

// Observer определяет интерфейс наблюдателя в паттерне Observer.
// Реализации этого интерфейса получают уведомления о событиях аудита.
type Observer interface {
	// Notify вызывается при возникновении события аудита
	Notify(event AuditEvent)
}

// Subject определяет интерфейс субъекта в паттерне Observer.
// Управляет списком наблюдателей и уведомляет их о событиях.
type Subject interface {
	// Attach добавляет наблюдателя в список
	Attach(observer Observer)

	// Detach удаляет наблюдателя из списка
	Detach(observer Observer)

	// NotifyAll уведомляет всех наблюдателей о событии
	NotifyAll(event AuditEvent)
}

// AuditService реализует сервис аудита и интерфейс Subject.
// Управляет наблюдателями и распространяет события аудита.
type AuditService struct {
	observers []Observer
	mu        sync.RWMutex
}

// NewAuditService создает новый экземпляр сервиса аудита.
//
// Возвращает:
//   - *AuditService: инициализированный сервис аудита
func NewAuditService() *AuditService {
	return &AuditService{
		observers: make([]Observer, 0),
	}
}

// Attach добавляет наблюдателя в список для получения уведомлений.
//
// Параметры:
//   - observer: наблюдатель для добавления
func (s *AuditService) Attach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

// Detach удаляет наблюдателя из списка уведомлений.
//
// Параметры:
//   - observer: наблюдатель для удаления
func (s *AuditService) Detach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// NotifyAll уведомляет всех зарегистрированных наблюдателей о событии.
// Уведомления отправляются асинхронно, чтобы не блокировать основной поток.
//
// Параметры:
//   - event: событие аудита для отправки наблюдателям
func (s *AuditService) NotifyAll(event AuditEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, observer := range s.observers {
		// Запускаем уведомления асинхронно, чтобы не блокировать основной поток
		go observer.Notify(event)
	}
}

// LogEvent создает и регистрирует событие аудита.
// Метод создает объект AuditEvent с текущим timestamp и уведомляет всех наблюдателей.
//
// Параметры:
//   - action: тип действия (ActionShorten или ActionFollow)
//   - userID: идентификатор пользователя
//   - url: оригинальный URL
func (s *AuditService) LogEvent(action, userID, url string) {
	event := AuditEvent{
		Timestamp: time.Now().Unix(),
		Action:    action,
		UserID:    userID,
		URL:       url,
	}
	s.NotifyAll(event)
}

// FileObserver наблюдатель для записи в файл
type FileObserver struct {
	filePath string
	mu       sync.Mutex
}

// NewFileObserver создает новый наблюдатель для записи в файл
func NewFileObserver(filePath string) *FileObserver {
	return &FileObserver{
		filePath: filePath,
	}
}

// Notify реализует интерфейс Observer
func (f *FileObserver) Notify(event AuditEvent) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Открываем файл для добавления
	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening audit file: %v", err)
		return
	}
	defer file.Close()

	// Сериализуем событие в JSON
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling audit event: %v", err)
		return
	}

	// Добавляем новую строку
	data = append(data, '\n')

	// Записываем в файл
	if _, err := file.Write(data); err != nil {
		log.Printf("Error writing to audit file: %v", err)
	}
}

// RemoteObserver наблюдатель для отправки на удаленный сервер
type RemoteObserver struct {
	url    string
	client *http.Client
}

// NewRemoteObserver создает новый наблюдатель для отправки на удаленный сервер
func NewRemoteObserver(url string) *RemoteObserver {
	return &RemoteObserver{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Notify реализует интерфейс Observer
func (r *RemoteObserver) Notify(event AuditEvent) {
	// Сериализуем событие в JSON
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling audit event for remote: %v", err)
		return
	}

	// Отправляем POST запрос
	resp, err := r.client.Post(r.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error sending audit event to remote server: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Remote audit server returned non-OK status: %d", resp.StatusCode)
	}
}
