package cntxt

import (
	"context"
	"log/slog"
	"opg-reports/report/package/logger"
	"os"
)

const loggerKey string = "logger"

// AddLogger
func AddLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

// WithLogger
func WithLogger(ctx context.Context) context.Context {
	var log *slog.Logger = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))
	return context.WithValue(ctx, loggerKey, log)
}

// GetLogger
func GetLogger(ctx context.Context) (log *slog.Logger) {
	if v := ctx.Value(loggerKey); v != nil {
		log = v.(*slog.Logger)
	} else {
		log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))
	}
	return

}
