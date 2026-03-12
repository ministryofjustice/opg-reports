// Package `ctxlog` expands on context and combines with slog
// to provide a direct way of getting and genrating a context
// attached logger
package ctxlog

import (
	"context"
	"log/slog"
	"opg-reports/report/packages/types"
)

// ctxLogger is a struct o handler extending the slog interface to
// directly fetch itself from context.Context.
//
// This is used heavily by all main functions to capture the current
// context.
//
// types.Logger
// types.ContextLogger
type ctxLogger struct {
	context.Context
	log types.Logger
}

// Log returns the logger instance
func (self *ctxLogger) Logger() types.Logger {
	return self.log
}

// Log returns the logger instance
func (self *ctxLogger) Log() *slog.Logger {
	return self.log.Log()
}

// Ctx returns the context
func (self *ctxLogger) Ctx() context.Context {
	return self
}

// New returns a new context logger
func New(ctx context.Context, logger types.Logger) types.ContextLogger {
	var lg types.Logger = logger
	var ctxLog *ctxLogger

	v := ctx.Value("LOGGER")
	// use the passed version if set
	// or if the passed logger is nol, and so is the one from context, create a new one
	// or logger is ni but we have a context version - use that
	if logger != nil {
		lg = logger
	} else if logger == nil && v == nil {
		lg = Log(nil, NIL)
		ctx = context.WithValue(ctx, "LOGGER", lg)
	} else if logger == nil && v != nil {
		lg = v.(*log)
	}
	// seupt and return
	ctxLog = &ctxLogger{
		Context: ctx,
		log:     lg,
	}
	return ctxLog

}
