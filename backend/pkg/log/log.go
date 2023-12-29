package log

import (
	"io"
	"log"
	"os"
)

// Service logger (now zap)
var loggerEntry *log.Logger

// Logger wrapper for logrus
type Logger struct {
	*log.Logger
}

func GetLogger() *Logger {
	return &Logger{loggerEntry}
}

// Init need for creating logger (now logrus)
func Init() {
	if loggerEntry != nil {
		return
	}

	err := os.MkdirAll("logs", 0755)
	if err != nil || os.IsExist(err) {
		panic(err)
	}

	logFile, err := os.OpenFile("logs/service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}

	logger := log.New(io.MultiWriter(logFile, os.Stderr), "LOG ", log.Ldate|log.Ltime|log.Llongfile)
	loggerEntry = logger
}
