package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[FORUM-DIAPASON] ", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string) {
	l.Println("INFO:", msg)
}

func (l *Logger) Error(msg string) {
	l.Println("ERROR:", msg)
}

func (l *Logger) Warn(msg string) {
	l.Println("WARN:", msg)
}
