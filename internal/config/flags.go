package config

import (
	"flag"
	"strings"
)

// parseFlags обрабатывает аргументы командной строки
func parseFlags() (string, string) {
	// адрес запуска HTTP-сервера значением :8080 по умолчанию
	portFlag := flag.String("a", "8080", "address and port to run server")

	// базовый адрес результирующего сокращённого URL значением
	resAddressFlag := flag.String("b", "http://localhost:8080", "address and port for short url")

	flag.Parse()

	port := strings.Split(*portFlag, ":")[1]
	resAddress := *resAddressFlag

	return port, resAddress
}
