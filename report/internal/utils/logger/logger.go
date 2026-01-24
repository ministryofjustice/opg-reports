package logger

import (
	"log/slog"
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
	// case "info":
	// 	return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

// New returns a configured slog.Logger instance
// that sets the log level and log handler.
//
// By default, the level is set to Info and TextHandler
func New(lvl string, as ...string) (logger *slog.Logger) {
	var (
		options *slog.HandlerOptions = &slog.HandlerOptions{Level: strToLevel(lvl)}
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
