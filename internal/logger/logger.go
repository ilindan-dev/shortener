package logger

import (
	"github.com/ilindan-dev/shortener/internal/config"
	"github.com/rs/zerolog"
	"os"
)

// NewLogger creates a new configured instance of zerolog.Logger.
func NewLogger(cfg *config.Config) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		// Default to info level if config is invalid or missing
		level = zerolog.InfoLevel
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05Z07:00"}

	logger := zerolog.New(consoleWriter).With().
		Timestamp().
		Str("service", "delayed-notifier").
		Caller().
		Logger().
		Level(level)

	return &logger, nil
}
