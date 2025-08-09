package config

// Структура для конфига
type ConfigStruct struct {
	Protocol     string
	Port         string
	ShortAddress string
}

// Глобальная переменная для конфига
var Configuration *ConfigStruct = GenerateConfig()

// Генерация конфигурации
func GenerateConfig() *ConfigStruct {
	// Получение данных из флагов
	reqAddr, resAddr := parseFlags()

	return &ConfigStruct{
		Protocol:     "http://",
		Port:         reqAddr,
		ShortAddress: resAddr,
	}
}
