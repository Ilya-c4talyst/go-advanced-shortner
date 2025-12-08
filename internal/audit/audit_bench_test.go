package audit

import (
	"os"
	"testing"
	"time"
)

// BenchmarkAuditService_LogEvent измеряет производительность логирования событий
func BenchmarkAuditService_LogEvent(b *testing.B) {
	service := NewAuditService()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		service.LogEvent(ActionShorten, "user123", "https://example.com/test")
	}
}

// BenchmarkFileObserver_Notify измеряет производительность записи в файл
func BenchmarkFileObserver_Notify(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "audit-bench-*.json")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	observer := NewFileObserver(tmpFile.Name())
	event := AuditEvent{
		Timestamp: time.Now().Unix(),
		Action:    ActionShorten,
		UserID:    "user123",
		URL:       "https://example.com/test",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		observer.Notify(event)
	}
}

// BenchmarkAuditService_WithMultipleObservers измеряет производительность с несколькими наблюдателями
func BenchmarkAuditService_WithMultipleObservers(b *testing.B) {
	service := NewAuditService()

	// Создаем временный файл для FileObserver
	tmpFile, err := os.CreateTemp("", "audit-bench-*.json")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	fileObserver := NewFileObserver(tmpFile.Name())
	service.Attach(fileObserver)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		service.LogEvent(ActionShorten, "user123", "https://example.com/test")
	}

	// Даем время для завершения асинхронных операций
	time.Sleep(100 * time.Millisecond)
}
