package repository

import (
	"fmt"
	"testing"
)

// BenchmarkMemoryRepository_SetValue измеряет производительность сохранения URL в памяти
func BenchmarkMemoryRepository_SetValue(b *testing.B) {
	repo := NewMemoryRepository()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		shortURL := fmt.Sprintf("short%d", i)
		originalURL := fmt.Sprintf("https://example.com/path/%d", i)
		_ = repo.SetValue(shortURL, originalURL, "user123")
	}
}

// BenchmarkMemoryRepository_GetFullValue измеряет производительность получения URL из памяти
func BenchmarkMemoryRepository_GetFullValue(b *testing.B) {
	repo := NewMemoryRepository()

	// Подготовка данных
	_ = repo.SetValue("test123", "https://example.com/test", "user123")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = repo.GetFullValue("test123")
	}
}

// BenchmarkMemoryRepository_GetUserURLs измеряет производительность получения URL пользователя
func BenchmarkMemoryRepository_GetUserURLs(b *testing.B) {
	repo := NewMemoryRepository()

	// Подготовка данных - создаем 100 URL для пользователя
	for i := 0; i < 100; i++ {
		shortURL := fmt.Sprintf("short%d", i)
		originalURL := fmt.Sprintf("https://example.com/path/%d", i)
		_ = repo.SetValue(shortURL, originalURL, "user123")
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = repo.GetUserURLs("user123")
	}
}

// BenchmarkMemoryRepository_SetValuesBatch измеряет производительность пакетного сохранения URL
func BenchmarkMemoryRepository_SetValuesBatch(b *testing.B) {
	repo := NewMemoryRepository()

	pairs := make(map[string]string, 100)
	for i := 0; i < 100; i++ {
		shortURL := fmt.Sprintf("short%d", i)
		originalURL := fmt.Sprintf("https://example.com/path/%d", i)
		pairs[shortURL] = originalURL
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = repo.SetValuesBatch(pairs, "user123")
	}
}

// BenchmarkMemoryRepository_DeleteURLsBatch измеряет производительность пакетного удаления URL
func BenchmarkMemoryRepository_DeleteURLsBatch(b *testing.B) {
	repo := NewMemoryRepository()

	// Подготовка данных
	shortURLs := make([]string, 10)
	for i := 0; i < 10; i++ {
		shortURL := fmt.Sprintf("short%d", i)
		originalURL := fmt.Sprintf("https://example.com/path/%d", i)
		_ = repo.SetValue(shortURL, originalURL, "user123")
		shortURLs[i] = shortURL
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = repo.DeleteURLsBatch(shortURLs, "user123")
	}
}
