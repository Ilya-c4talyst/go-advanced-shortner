package main

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/config"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/repository"
	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/service"
)

// TestMemoryProfile профилирует использование памяти
func TestMemoryProfile(t *testing.T) {
	cfg := &config.ConfigStruct{
		Protocol:     "http://",
		Port:         ":8080",
		ShortAddress: "http://localhost:8080",
		FilePath:     "",
		AddressDB:    "",
	}

	repo := repository.NewMemoryRepository()
	svc := service.NewURLShortnerService(repo, cfg)

	// Создаем большое количество URL для профилирования
	for i := 0; i < 10000; i++ {
		url := fmt.Sprintf("https://example.com/path/%d", i)
		userID := fmt.Sprintf("user%d", i%100)
		_, _ = svc.CreateShortURL(url, userID)
	}

	// Получаем URL пользователей
	for i := 0; i < 100; i++ {
		userID := fmt.Sprintf("user%d", i)
		_, _ = svc.GetUserURLs(userID)
	}

	// Принудительный сбор мусора для точности профилирования
	runtime.GC()
}

// BenchmarkMemoryProfile бенчмарк для профилирования памяти
func BenchmarkMemoryProfile(b *testing.B) {
	cfg := &config.ConfigStruct{
		Protocol:     "http://",
		Port:         ":8080",
		ShortAddress: "http://localhost:8080",
		FilePath:     "",
		AddressDB:    "",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		repo := repository.NewMemoryRepository()
		svc := service.NewURLShortnerService(repo, cfg)

		// Создаем URL
		for j := 0; j < 1000; j++ {
			url := fmt.Sprintf("https://example.com/path/%d", j)
			userID := fmt.Sprintf("user%d", j%10)
			_, _ = svc.CreateShortURL(url, userID)
		}

		// Получаем URL пользователей
		for j := 0; j < 10; j++ {
			userID := fmt.Sprintf("user%d", j)
			_, _ = svc.GetUserURLs(userID)
		}
	}
}
