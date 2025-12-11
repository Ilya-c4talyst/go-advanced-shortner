package config

// ConfigStruct хранит конфигурацию приложения.
// Структура содержит все настройки, необходимые для работы сервиса сокращения URL,
// включая сетевые параметры, пути к файлам, настройки базы данных и аудита.
type ConfigStruct struct {
	// Protocol определяет протокол для HTTP-сервера (обычно "http://")
	Protocol string

	// Port определяет порт, на котором будет запущен HTTP-сервер (например, ":8080")
	Port string

	// ShortAddress определяет базовый адрес для коротких URL (например, "http://localhost:8080")
	ShortAddress string

	// FilePath определяет путь к файлу для хранения URL (если используется файловое хранилище)
	FilePath string

	// AddressDB определяет строку подключения к базе данных PostgreSQL
	AddressDB string

	// AuthSecretKey определяет секретный ключ для подписи токенов аутентификации
	AuthSecretKey string

	// AuditFile определяет путь к файлу для записи событий аудита
	AuditFile string

	// AuditURL определяет URL удаленного сервера для отправки событий аудита
	AuditURL string
}

// GenerateConfig создает и возвращает конфигурацию приложения.
// Функция парсит флаги командной строки и переменные окружения,
// создавая полностью инициализированную структуру конфигурации.
//
// Возвращает:
//   - *ConfigStruct: заполненная структура конфигурации
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
