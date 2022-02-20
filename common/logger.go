package common

import (
	"log"
	"os"
)

type Logger struct {
	logFile *os.File
	logger  *log.Logger
}

func NewLogger(name string) *Logger {
	logFile, _ := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0744)
	logger := log.New(logFile, "INFO ", log.Ldate|log.Ltime)

	return &Logger{logFile, logger}
}

func (logger *Logger) Log(messages ...interface{}) {
	logger.logger.Println(messages...)
}
