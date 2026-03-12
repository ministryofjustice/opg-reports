package ctxlog

import (
	"log/slog"
	"opg-reports/report/packages/types"
	"os"
	"strings"
)

type handlerType string

const (
	NIL  handlerType = ""
	TEXT handlerType = "TEXT"
	JSON handlerType = "JSON"
)

// string to level
var logLevels = map[string]slog.Leveler{
	"ERROR": slog.LevelError,
	"WARN":  slog.LevelWarn,
	"INFO":  slog.LevelWarn,
	"DEBUG": slog.LevelDebug,
}

// log is an internal struct to deal with the logging
// and wrapping around slog
type log struct {
	logger  *slog.Logger
	leveler slog.Leveler
	handler slog.Handler
}

// Log returns the logger itself
func (self *log) Log() *slog.Logger {
	return self.logger
}

// Leveler returns the log level
func (self *log) Leveler() slog.Leveler {
	return self.leveler
}

// Handler returns the log handler - text/json
func (self *log) Handler() slog.Handler {
	return self.handler
}

// leveler is internal to create a leveler from either the
// value passed or the env, defaulting to info
func leveler(l slog.Leveler) slog.Leveler {
	v := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	// map the env log level to a real one
	if lv, ok := logLevels[v]; l == nil && ok {
		l = lv
	}
	// fallback to info
	if l == nil {
		l = slog.LevelInfo
	}

	return l
}

// handler configures the handler based off passed values or the env;
// defaults to text
func handler(l slog.Leveler, ht handlerType) slog.Handler {
	var h slog.Handler
	var v = strings.ToUpper(os.Getenv("LOG_TYPE"))

	if ht == JSON {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l})
	} else if ht == TEXT {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: l})
	} else if v == string(JSON) {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: l})
	}

	return h
}

// Log creates a standalone logger instance
func Log(l slog.Leveler, ht handlerType) types.Logger {
	var h slog.Handler
	l = leveler(l)
	h = handler(l, ht)
	return &log{
		leveler: l,
		handler: h,
		logger:  slog.New(h),
	}
}
