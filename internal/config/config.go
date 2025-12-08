package config

// Структура для конфига
type ConfigStruct struct {
	Protocol      string
	Port          string
	ShortAddress  string
	FilePath      string
	AddressDB     string
	AuthSecretKey string
	AuditFile     string // путь к файлу для аудита
	AuditURL      string // URL для удаленного сервера аудита
}

// Генерация конфигурации
func GenerateConfig() *ConfigStruct {
	// Получение данных из флагов
	reqAddr, resAddr, filePath, dbAddress, auditFile, auditURL := parseFlags()

	return &ConfigStruct{
		Protocol:      "http://",
		Port:          reqAddr,
		ShortAddress:  resAddr,
		FilePath:      filePath,
		AddressDB:     dbAddress,
		AuthSecretKey: "your-secret-key-change-in-production", // В продакшене должен быть из переменной окружения
		AuditFile:     auditFile,
		AuditURL:      auditURL,
	}
}
