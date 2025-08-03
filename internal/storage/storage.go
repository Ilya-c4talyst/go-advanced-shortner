package storage

// Тип для БД
type DB map[string]string

// Создание БД (пока мапа, потом поправим)
func CreateDB() DB {
	return make(map[string]string)
}
