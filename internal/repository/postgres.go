package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// PostgreSQLRepository реализация репозитория для работы с PostgreSQL
type PostgreSQLRepository struct {
	pool *pgxpool.Pool
}

// NewPostgreSQLRepository создает новый репозиторий для работы с PostgreSQL
func NewPostgreSQLRepository(dsn string) (URLRepository, error) {
	// Подключаемся к базе данных
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	repo := &PostgreSQLRepository{
		pool: pool,
	}

	// Выполняем миграции
	if err := repo.runMigrations(dsn); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return repo, nil
}

// runMigrations выполняет миграции базы данных
func (r *PostgreSQLRepository) runMigrations(dsn string) error {
	// Открываем соединение через database/sql для миграций
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// GetValue получает оригинальный URL по короткому
func (r *PostgreSQLRepository) GetValue(shortURL string) (string, error) {
	var originalURL string
	err := r.pool.QueryRow(context.Background(), 
		"SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("not found key in database")
		}
		return "", fmt.Errorf("failed to get value: %v", err)
	}
	
	return originalURL, nil
}

// SetValue сохраняет пару короткий URL - оригинальный URL
func (r *PostgreSQLRepository) SetValue(shortURL, originalURL string) error {
	_, err := r.pool.Exec(context.Background(),
		"INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO UPDATE SET original_url = $2",
		shortURL, originalURL)
	
	if err != nil {
		return fmt.Errorf("failed to set value: %v", err)
	}
	
	return nil
}

// Close закрывает соединение с базой данных
func (r *PostgreSQLRepository) Close() error {
	r.pool.Close()
	return nil
}
