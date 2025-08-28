package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Ilya-c4talyst/go-advanced-shortner/internal/model"
)

// Структура для БД
type DB struct {
	data    map[string]string
	counter int
}

// Создание БД
func CreateDB() *DB {
	return &DB{
		data:    make(map[string]string),
		counter: 0,
	}
}

// Получение значения по ключу
func (db *DB) Get(key string) (string, bool) {
	value, exists := db.data[key]
	return value, exists
}

// Установка значения
func (db *DB) Set(key, value string) {
	db.data[key] = value
	db.counter++
}

// Сохраняет данные в файл в формате JSON
func (db *DB) SaveToFile(filePath string) error {
	var records []model.URLRecord
	
	counter := 1
	for shortURL, originalURL := range db.data {
		record := model.URLRecord{
			UUID:        strconv.Itoa(counter),
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		}
		records = append(records, record)
		counter++
	}
	
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	
	// Создаем директорию если её нет
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(filePath, data, 0644)
}

// Загружает данные из файла в формате JSON
func (db *DB) LoadFromFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // Если файл не существует, просто возвращаем nil
	}
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	// Пустой файл
	if len(data) == 0 {
		return nil
	}
	
	var records []model.URLRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}
	
	// Загружаем данные в мапу
	maxUUID := 0
	for _, record := range records {
		db.data[record.ShortURL] = record.OriginalURL
		if uuid, err := strconv.Atoi(record.UUID); err == nil && uuid > maxUUID {
			maxUUID = uuid
		}
	}
	db.counter = maxUUID
	
	return nil
}
