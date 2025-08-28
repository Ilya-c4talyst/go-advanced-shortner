package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

// parseFlags обрабатывает аргументы командной строки
func parseFlags() (string, string, string) {
	// адрес запуска HTTP-сервера значением localhost:8080 по умолчанию
	portFlag := flag.String("a", "localhost:8080", "address and port to run server")

	// базовый адрес результирующего сокращённого URL значением
	resAddressFlag := flag.String("b", "http://localhost:8080", "address and port for short url")

	// Добавление нового флага и переменной окружения для пути к файлу
	filePathFlag := flag.String("f", "data/storage.json", "path to the file for storing data")

	flag.Parse()

	// Проверка и обработка адреса сервера
	portParts := strings.Split(*portFlag, ":")
	if len(portParts) < 2 {
		log.Printf("invalid address format: %s, expected format: host:port\n", *portFlag)
		flag.Usage()
		os.Exit(2)
	}

	port := ":" + portParts[1]

	// Проверка базового адреса
	resAddress := *resAddressFlag
	if !strings.HasPrefix(resAddress, "http://") && !strings.HasPrefix(resAddress, "https://") {
		log.Printf("invalid base address: %s, must start with http:// or https://\n", resAddress)
		flag.Usage()
		os.Exit(2)
	}

	// Путь до файла с urls
	filePath := *filePathFlag

	// Если параметры заданы через переменные окружения, используем их
	if os.Getenv("SERVER_ADDRESS") != "" {
		port = os.Getenv("SERVER_ADDRESS")
	}
	if os.Getenv("BASE_URL") != "" {
		resAddress = os.Getenv("BASE_URL")
	}
	if os.Getenv("FILE_STORAGE_PATH") != "" {
		filePath = os.Getenv("FILE_STORAGE_PATH")
	}

	return port, resAddress, filePath
}
