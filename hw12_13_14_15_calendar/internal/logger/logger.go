package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levels = map[string]zapcore.Level{
	"INFO":  zapcore.InfoLevel,
	"DEBUG": zapcore.DebugLevel,
	"ERROR": zapcore.ErrorLevel,
}

type Logger = zap.SugaredLogger

func New(loggerConf *config.LoggerConf) *Logger {
	curDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("can't get working dir"))
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logPath := filepath.Join(curDir, loggerConf.Path)
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(fmt.Errorf("can't open log file"))
	}
	writer := zapcore.AddSync(logFile)

	lvl, ok := levels[loggerConf.Level]
	var defaultLogLevel zapcore.Level
	if !ok {
		fmt.Printf("no such logger level: %v", lvl)
		defaultLogLevel = zapcore.DebugLevel
	} else {
		defaultLogLevel = lvl
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	sugar := logger.Sugar()
	return sugar
}
