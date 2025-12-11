package audit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestAuditService_AttachDetach(t *testing.T) {
	service := NewAuditService()

	// Создаем мок-наблюдателя
	observer := &mockObserver{}

	// Проверяем добавление наблюдателя
	service.Attach(observer)
	if len(service.observers) != 1 {
		t.Errorf("Expected 1 observer, got %d", len(service.observers))
	}

	// Проверяем удаление наблюдателя
	service.Detach(observer)
	if len(service.observers) != 0 {
		t.Errorf("Expected 0 observers, got %d", len(service.observers))
	}
}

func TestAuditService_LogEvent(t *testing.T) {
	service := NewAuditService()
	observer := &mockObserver{
		events: make([]AuditEvent, 0),
	}

	service.Attach(observer)

	// Логируем событие
	service.LogEvent(ActionShorten, "user123", "https://example.com")

	// Даем время для асинхронной обработки
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что событие было получено
	if len(observer.events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(observer.events))
	}

	if observer.events[0].Action != ActionShorten {
		t.Errorf("Expected action %s, got %s", ActionShorten, observer.events[0].Action)
	}

	if observer.events[0].UserID != "user123" {
		t.Errorf("Expected user_id user123, got %s", observer.events[0].UserID)
	}

	if observer.events[0].URL != "https://example.com" {
		t.Errorf("Expected url https://example.com, got %s", observer.events[0].URL)
	}
}

func TestFileObserver_Notify(t *testing.T) {
	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", "audit-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	observer := NewFileObserver(tmpFile.Name())

	event := AuditEvent{
		Timestamp: time.Now().Unix(),
		Action:    ActionShorten,
		UserID:    "user123",
		URL:       "https://example.com",
	}

	// Отправляем событие
	observer.Notify(event)

	// Читаем файл
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read audit file: %v", err)
	}

	// Проверяем, что событие записано
	var savedEvent AuditEvent
	if err := json.Unmarshal(data, &savedEvent); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	if savedEvent.Action != event.Action {
		t.Errorf("Expected action %s, got %s", event.Action, savedEvent.Action)
	}

	if savedEvent.UserID != event.UserID {
		t.Errorf("Expected user_id %s, got %s", event.UserID, savedEvent.UserID)
	}

	if savedEvent.URL != event.URL {
		t.Errorf("Expected url %s, got %s", event.URL, savedEvent.URL)
	}
}

func TestRemoteObserver_Notify(t *testing.T) {
	// Создаем мок-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		var event AuditEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			t.Errorf("Failed to decode event: %v", err)
		}

		if event.Action != ActionFollow {
			t.Errorf("Expected action %s, got %s", ActionFollow, event.Action)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	observer := NewRemoteObserver(server.URL)

	event := AuditEvent{
		Timestamp: time.Now().Unix(),
		Action:    ActionFollow,
		UserID:    "user456",
		URL:       "https://example.com/long/path",
	}

	// Отправляем событие
	observer.Notify(event)

	// Даем время для отправки
	time.Sleep(100 * time.Millisecond)
}

// mockObserver для тестов
type mockObserver struct {
	events []AuditEvent
}

func (m *mockObserver) Notify(event AuditEvent) {
	m.events = append(m.events, event)
}
