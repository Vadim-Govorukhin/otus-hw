package logger

import "fmt"

type Logger struct { // TODO
	Level string // Уровень логирования
}

func New(level string) *Logger {
	return &Logger{Level: level}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

// TODO
