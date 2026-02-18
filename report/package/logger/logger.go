package logger

import (
	"log/slog"
	"opg-reports/report/package/env"
	"os"
	"strings"
)

func strToLevel(lvl string) slog.Leveler {
	lvl = strings.ToLower(lvl)
	switch lvl {
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "debug":
		return slog.LevelDebug
	}
	// default to info
	return slog.LevelInfo
}

// New returns a configured slog.Logger instance
// that sets the log level and log handler.
//
// # By default, the level is set to Info and TextHandler
//
// Log level is overwritten from environment variables
func New(lvl string, as ...string) (logger *slog.Logger) {
	var (
		level   string               = env.Get("LOG_LEVEL", lvl)
		options *slog.HandlerOptions = &slog.HandlerOptions{Level: strToLevel(level)}
		handler slog.Handler         = slog.NewTextHandler(os.Stdout, options)
	)

	// hackey way to have optional param
	for _, h := range as {
		h = strings.ToLower(h)
		if h == "json" {
			logger = slog.New(slog.NewJSONHandler(os.Stdout, options))
			break
		}
	}
	logger = slog.New(handler)
	return

}
