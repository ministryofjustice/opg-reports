package logx

import (
	"log/slog"
	"opg-reports/app/internal/envx"
	"os"
	"strings"
)

// Default is a wrapper around slog.Default
//
// Use this instead so the logx init will have run and
// therefore pulled in defaults from the env
func Default() *slog.Logger {
	return slog.Default()
}

// Set configures the default slog instance with standard
// approach (json, source added) and allows configuration
// of the log level
func Set(l slog.Leveler) *slog.Logger {
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

	return slog.Default()
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

func init() {
	Set(level())
}
