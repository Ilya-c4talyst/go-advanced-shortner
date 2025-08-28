package storage

import (
	"encoding/json"
	"os"
)

// Тип для БД
type DB map[string]string

// Создание БД (пока мапа, потом поправим)
func CreateDB() DB {
	return make(map[string]string)
}

// Сохраняет данные в файл в формате JSON
func (db DB) SaveToFile(filePath string) error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// Загружает данные из файла в формате JSON
func (db DB) LoadFromFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // Если файл не существует, просто возвращаем nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &db)
}
