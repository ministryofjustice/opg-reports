package logx

import (
	"log/slog"
	"opg-reports/app/internal/envx"
	"os"
	"strings"
)

// Default sets the default logging to be used by all apps.
// Sets the log level from environment variable (`LOG_LEVEL`)
// and defaults to info.
//
// Always uses JSON handler & adds the source details
func Default() {
	Set(level())
}

// Set configures the default slog instance with standard
// approach (json, source added) and allows configuration
// of the log level
func Set(l slog.Leveler) {
	var handler slog.Handler
	var logger *slog.Logger

	handler = slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     l,
		},
	)
	logger = slog.New(handler)

	slog.SetDefault(logger)
}

// level grab log level from env
func level() (l slog.Leveler) {
	var key = "LOG_LEVEL"
	var envValue = envx.Get(key, "info")

	switch strings.ToLower(envValue) {
	case "error", "err", "e":
		l = slog.LevelError
	case "warning", "warn", "w":
		l = slog.LevelWarn
	case "debugging", "debug", "d":
		l = slog.LevelDebug
	default:
		l = slog.LevelInfo
	}

	return
}
