// Package logger providers methods to create an slog instance with context.
//
// Log level & handler type can be configured via environment variables or
// from context values - with fallbacks if not set.
//
// Will pull existing logger from context where it can, otherwise creates
// and attaches a new one.
package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

const (
	ctxKey        string = "logger"
	ctxLevelKey   string = "LOG_LEVEL"
	ctxHandlerKey string = "LOG_HANDLER"
)

// Level returns the log level based in priority order...
//
// - from environment variable `LOG_LEVEL“
// - from context value `LOG_LEVEL`
// - from the non nil `lvl` parameter
// - fallback fixed value of `info`
func Level(ctx context.Context, lvl *string) (c context.Context, l slog.Leveler) {
	var str string = ""
	c = ctx
	// get the level name as a string
	if v := os.Getenv(ctxLevelKey); v != "" {
		str = v
	} else if v := ctx.Value(ctxLevelKey); v != nil && v.(string) != "" {
		str = v.(string)
	} else if lvl != nil {
		str = *lvl
	}
	str = strings.ToLower(str)
	// now get the leveler from the string value
	switch str {
	case "error", "err", "e":
		l = slog.LevelError
	case "warn", "w":
		l = slog.LevelWarn
	case "debug", "d":
		l = slog.LevelDebug
	default:
		l = slog.LevelInfo
	}
	c = context.WithValue(ctx, ctxLevelKey, str)

	return
}

// Handler returns the log handler based on settings in this order:
//
// - from environment variable `LOG_HANDLER
// - from context value `LOG_HANDLER`
// - from the non nil `logType` parameter
// - fallback to `NewTextHandler`
func Handler(ctx context.Context, logType *string, options *slog.HandlerOptions) (c context.Context, h slog.Handler) {
	var str string = ""
	c = ctx
	// get the level name as a string
	if v := os.Getenv(ctxHandlerKey); v != "" {
		str = v
	} else if v := ctx.Value(ctxHandlerKey); v != nil && v.(string) != "" {
		str = v.(string)
	} else if logType != nil {
		str = *logType
	}
	// check value
	str = strings.ToLower(str)
	switch str {
	case "json":
		h = slog.NewJSONHandler(os.Stdout, options)
	default:
		h = slog.NewTextHandler(os.Stdout, options)
	}
	c = context.WithValue(ctx, ctxHandlerKey, str)
	return
}

// New create a brand new logger instance with values set from env / ctx / inputs
func New(ctx context.Context, lvl *string, logType *string) (c context.Context, logger *slog.Logger) {
	var (
		level   slog.Leveler
		handler slog.Handler
	)
	c, level = Level(ctx, lvl)
	c, handler = Handler(ctx, logType, &slog.HandlerOptions{Level: level})

	logger = slog.New(handler)

	c = context.WithValue(ctx, ctxKey, logger)
	return
}

// Get returns an existing logger from the context where possible.
//
// If none is found then a new logger is created that will use the
// env & context values for the level & type
func Get(ctx context.Context) (c context.Context, logger *slog.Logger) {
	c = ctx
	if v := ctx.Value(ctxKey); v != nil {
		logger = v.(*slog.Logger)
	} else {
		c, logger = New(ctx, nil, nil)
	}
	return
}
