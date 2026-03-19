package slogx

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// HandlerType
type HandlerType string

const (
	JSON HandlerType = `JSON`
	TEXT HandlerType = `TEXT`
)

const (
	ctxLogKey        string = `logger`
	envLogLevelKey   string = `LOG_LEVEL`
	envLogHandlerKey string = `LOG_TYPE`
)

// Logger interface struct for exposing the core
// methods we use within the app code
type Logger interface {
	// Leveler returns the current leveler - used to check log levels etc
	Leveler() slog.Leveler
	// LevelStr is a helper function to return the log level as a string directly
	LevelStr() string
	// Handler returns the current handler
	Handler() slog.Handler
	// HandlerStr returns the handler type as a string (JSON/TEXT)
	HandlerStr() string
	// Debug forces ctx usage version
	Debug(ctx context.Context, msg string, args ...any)
	// Info forces ctx usage version
	Info(ctx context.Context, msg string, args ...any)
	// Warn forces ctx usage version
	Warn(ctx context.Context, msg string, args ...any)
	// Error forces ctx usage version
	Error(ctx context.Context, msg string, args ...any)
}

type logger struct {
	*slog.Logger
	leveler     slog.Leveler
	handler     slog.Handler
	handlerType HandlerType
}

// Leveler returns the current leveler - used to check log levels etc
func (self *logger) Leveler() slog.Leveler {
	return self.leveler
}

// LevelStr is a helper function to return the log level as a string directly
func (self *logger) LevelStr() string {
	return self.leveler.Level().String()
}

// Handler returns the current handler
func (self *logger) Handler() slog.Handler {
	return self.handler
}

// HandlerStr returns the handler type as a string (JSON/TEXT)
func (self *logger) HandlerStr() string {
	return string(self.handlerType)
}

// Debug forces ctx usage version
func (self *logger) Debug(ctx context.Context, msg string, args ...any) {
	self.Logger.DebugContext(ctx, msg, args...)
}

// Info forces ctx usage version
func (self *logger) Info(ctx context.Context, msg string, args ...any) {
	self.Logger.InfoContext(ctx, msg, args...)
}

// Warn forces ctx usage version
func (self *logger) Warn(ctx context.Context, msg string, args ...any) {
	self.Logger.WarnContext(ctx, msg, args...)
}

// Error forces ctx usage version
func (self *logger) Error(ctx context.Context, msg string, args ...any) {
	self.Logger.ErrorContext(ctx, msg, args...)
}

// New creates a new logger
func New(leveler slog.Leveler, handlerType HandlerType) (log Logger) {
	var handler slog.Handler
	var options = &slog.HandlerOptions{Level: leveler}

	switch handlerType {
	case JSON:
		handler = slog.NewJSONHandler(os.Stdout, options)
	default:
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	log = &logger{
		leveler:     leveler,
		handler:     handler,
		handlerType: handlerType,
		Logger:      slog.New(handler),
	}

	return

}

// FromContext pulls logger from an existing context, if it cant find one
// it will create a new logger with default values (info & text)
func FromContext(ctx context.Context) (log Logger) {
	if v := ctx.Value(ctxLogKey); v != nil {
		log = v.(Logger)
	} else {
		log = New(slog.LevelWarn, TEXT)
	}
	return
}

// Attach attachs the logger to the context via a known key
func Attach(ctx context.Context, log Logger) context.Context {
	ctx = context.WithValue(ctx, ctxLogKey, log)
	return ctx
}

// Config returns the constructors for New
func Config() (lvl slog.Leveler, h HandlerType) {
	lvl = leveler()
	h = handlerType()
	return
}

// handlerType returns the HandlerType to use checking the env
func handlerType() (h HandlerType) {
	var osHandler string = strings.ToUpper(os.Getenv(envLogHandlerKey))
	h = TEXT
	if osHandler == "" {
		return
	}
	if osHandler == string(JSON) {
		h = JSON
	}
	return
}

// leveler returns a leveler, checking values in the env
func leveler() (lvl slog.Leveler) {
	var osLevel string = strings.ToLower(os.Getenv(envLogLevelKey))
	lvl = slog.LevelInfo
	if osLevel == "" {
		return
	}
	switch osLevel {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	return
}
