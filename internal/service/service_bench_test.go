package service

import (
	"fmt"
	"testing"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
)

// BenchmarkCreateShortURL измеряет производительность создания коротких URL
func BenchmarkCreateShortURL(b *testing.B) {
	cfg := &config.ConfigStruct{
		Protocol:     "http://",
		Port:         ":8080",
		ShortAddress: "http://localhost:8080",
		FilePath:     "",
		AddressDB:    "",
	}

	repo := repository.NewMemoryRepository()
	service := NewURLShortnerService(repo, cfg)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		url := fmt.Sprintf("https://example.com/path/%d", i)
		_, _ = service.CreateShortURL(url, "user123")
	}
}

// BenchmarkGetFullURL измеряет производительность получения полных URL
func BenchmarkGetFullURL(b *testing.B) {
	cfg := &config.ConfigStruct{
		Protocol:     "http://",
		Port:         ":8080",
		ShortAddress: "http://localhost:8080",
		FilePath:     "",
		AddressDB:    "",
	}

	repo := repository.NewMemoryRepository()
	service := NewURLShortnerService(repo, cfg)

	// Подготовка данных
	shortURL, _ := service.CreateShortURL("https://example.com/test", "user123")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = service.GetFullURL(shortURL)
	}
}

// BenchmarkCreateShortURLsBatch измеряет производительность пакетного создания URL
func BenchmarkCreateShortURLsBatch(b *testing.B) {
	cfg := &config.ConfigStruct{
		Protocol:     "http://",
		Port:         ":8080",
		ShortAddress: "http://localhost:8080",
		FilePath:     "",
		AddressDB:    "",
	}

	repo := repository.NewMemoryRepository()
	service := NewURLShortnerService(repo, cfg)

	urls := make([]string, 100)
	for i := 0; i < 100; i++ {
		urls[i] = fmt.Sprintf("https://example.com/path/%d", i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = service.CreateShortURLsBatch(urls, "user123")
	}
}

// BenchmarkGetUserURLs измеряет производительность получения URL пользователя
func BenchmarkGetUserURLs(b *testing.B) {
	cfg := &config.ConfigStruct{
		Protocol:     "http://",
		Port:         ":8080",
		ShortAddress: "http://localhost:8080",
		FilePath:     "",
		AddressDB:    "",
	}

	repo := repository.NewMemoryRepository()
	service := NewURLShortnerService(repo, cfg)

	// Подготовка данных - создаем несколько URL для пользователя
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("https://example.com/path/%d", i)
		_, _ = service.CreateShortURL(url, "user123")
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = service.GetUserURLs("user123")
	}
}
