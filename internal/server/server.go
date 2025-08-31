package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/storage"
)

// Представляет HTTP сервер с graceful shutdown
type Server struct {
	httpServer *http.Server
	db         *storage.DB
	filePath   string
}

// Создаёт новый сервер
func NewServer(addr string, handler http.Handler, db *storage.DB, filePath string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		db:       db,
		filePath: filePath,
	}
}

// Запускает сервер с graceful shutdown
func (s *Server) Start() error {
	// Канал для получения сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Сервер запущен на %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	<-quit
	log.Println("Завершение работы сервера...")

	return s.shutdown()
}

// graceful shutdown
func (s *Server) shutdown() error {
	// Создаём контекст с таймаутом для shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Завершаем HTTP сервер
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при завершении сервера: %v", err)
		return err
	}

	// Сохраняем данные перед завершением
	return s.saveData()
}

// Сохраняет данные в файл
func (s *Server) saveData() error {
	log.Println("Сохранение данных...")
	if err := s.db.SaveToFile(s.filePath); err != nil {
		log.Printf("Ошибка сохранения данных: %v", err)
		return err
	}
	
	log.Println("Данные успешно сохранены")
	return nil
}
