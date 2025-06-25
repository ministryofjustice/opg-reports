package utils

import (
	"log/slog"
	"os"
	"strings"
)

func Logger(lvl string, as string) (logger *slog.Logger) {
	var options = &slog.HandlerOptions{}

	switch lvl {
	case "ERROR", "error":
		options.Level = slog.LevelError
	case "WARN", "warn":
		options.Level = slog.LevelWarn
	case "INFO", "info":
		options.Level = slog.LevelInfo
	case "DEBUG", "debug":
		options.Level = slog.LevelDebug
	default:
		options.Level = slog.LevelInfo
	}

	as = strings.ToLower(as)
	if as == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, options))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, options))
	}
	return

}
