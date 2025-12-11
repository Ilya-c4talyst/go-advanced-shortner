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

// AuditEvent представляет событие аудита
type AuditEvent struct {
	Timestamp int64  `json:"ts"`      // unix timestamp события
	Action    string `json:"action"`  // действие: shorten или follow
	UserID    string `json:"user_id"` // идентификатор пользователя
	URL       string `json:"url"`     // оригинальный (не сокращенный) URL
}

// Action типы действий
const (
	ActionShorten = "shorten" // создание сокращенного URL
	ActionFollow  = "follow"  // прохождение по ссылке
)

// Observer интерфейс наблюдателя
type Observer interface {
	Notify(event AuditEvent)
}

// Subject интерфейс субъекта
type Subject interface {
	Attach(observer Observer)
	Detach(observer Observer)
	NotifyAll(event AuditEvent)
}

// AuditService сервис аудита (Subject)
type AuditService struct {
	observers []Observer
	mu        sync.RWMutex
}

// NewAuditService создает новый сервис аудита
func NewAuditService() *AuditService {
	return &AuditService{
		observers: make([]Observer, 0),
	}
}

// Attach добавляет наблюдателя
func (s *AuditService) Attach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

// Detach удаляет наблюдателя
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

// NotifyAll уведомляет всех наблюдателей
func (s *AuditService) NotifyAll(event AuditEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, observer := range s.observers {
		// Запускаем уведомления асинхронно, чтобы не блокировать основной поток
		go observer.Notify(event)
	}
}

// LogEvent логирует событие аудита
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
