package logger

import (
	"log"
	"os"
)

var logger *log.Logger

func GetLogger() *log.Logger {
	if logger == nil {
		return log.New(os.Stdout, "[PHP-LSP] ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return logger
}

func SetLogger(l *log.Logger) {
	logger = l
}

func CreateLogFile(logFile string) *log.Logger {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to open log file")
	}

	return log.New(file, "[PHP-LSP] ", log.Ldate|log.Ltime|log.Lshortfile)
}
