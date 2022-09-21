package logger

import (
	"fmt"
	"os"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
)

type Logger struct { // TODO
	Level string // Уровень логирования
}

func New(loggerConf *config.LoggerConf) *Logger {
	return &Logger{Level: loggerConf.Level}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	// TODO
}

// TODO
