package config

import (
	"flag"
)

// parseFlags обрабатывает аргументы командной строки
func parseFlags() (string, string) {

	// адрес запуска HTTP-сервера значением :8080 по умолчанию
	port := ":" + *flag.String("a", "8080", "address and port to run server")

	// базовый адрес результирующего сокращённого URL значением
	resAddress := *flag.String("b", "http://localhost:8080", "address and port for short url")

	flag.Parse()

	return port, resAddress
}
