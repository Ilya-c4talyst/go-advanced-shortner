package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

// parseFlags обрабатывает аргументы командной строки
func parseFlags() (string, string) {
	// адрес запуска HTTP-сервера значением localhost:8080 по умолчанию
	portFlag := flag.String("a", "localhost:8080", "address and port to run server")

	// базовый адрес результирующего сокращённого URL значением
	resAddressFlag := flag.String("b", "http://localhost:8080", "address and port for short url")

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

	return port, resAddress
}
