package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(level string) *Logger {
	logLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(logLevel)

	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	return &Logger{logger}
}
