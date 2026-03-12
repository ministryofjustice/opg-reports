package ctxlog

import (
	"fmt"
	"log/slog"
	"opg-reports/report/packages/types"
	"os"
	"testing"
)

var (
	_ types.Logger = &log{}
)

func TestPackagesCtxLogLogFromEnv(t *testing.T) {
	// test log level and handler from env variables
	os.Setenv("LOG_LEVEL", "WARN")
	os.Setenv("LOG_TYPE", "JSON")

	lg := Log(nil, NIL)
	// check the log level
	level := lg.Leveler().Level()
	if level != slog.LevelWarn {
		t.Errorf("mismatching log level")
	}
	// check the handler type
	handler := lg.Handler()
	if fmt.Sprintf("%T", handler) != "*slog.JSONHandler" {
		t.Errorf("mismatch handler: %T", handler)
	}

}

func TestPackagesCtxLogLogFromValues(t *testing.T) {

	lg := Log(slog.LevelError, JSON)
	// check the log level
	level := lg.Leveler().Level()
	if level != slog.LevelError {
		t.Errorf("mismatching log level")
	}
	// check the handler type
	handler := lg.Handler()
	if fmt.Sprintf("%T", handler) != "*slog.JSONHandler" {
		t.Errorf("mismatch handler: %T", handler)
	}

}

// func TestPackagesLoggerHandler(t *testing.T) {
// 	var log *Logger
// 	var handler slog.Handler
// 	var hType string
// 	var ctx context.Context = t.Context()
// 	var key string = "LOG_TYPE"

// 	// try the handler as json via argument
// 	hType = "json"
// 	log = logger(ctx)
// 	handler = log.Handler(&hType)
// 	if fmt.Sprintf("%T", handler) != "*slog.JSONHandler" {
// 		t.Errorf("incorrect handler returned")
// 	}
// 	if log.Value(key).(string) != "JSON" {
// 		t.Errorf("incorrect handler set in context")
// 	}

// 	// now set it via os
// 	hType = "text"
// 	os.Setenv(key, hType)
// 	log = logger(ctx)
// 	handler = log.Handler(nil)
// 	if fmt.Sprintf("%T", handler) != "*slog.TextHandler" {
// 		t.Errorf("incorrect handler returned: %s", fmt.Sprintf("%T", handler))
// 	}
// 	if log.Value(key).(string) != "TEXT" {
// 		t.Errorf("incorrect handler set in context")
// 	}

// }

// func TestPackagesLoggerLevel(t *testing.T) {
// 	var log *Logger
// 	var leveler slog.Leveler
// 	var lvl string
// 	var ctx context.Context = t.Context()
// 	var key string = "LOG_LEVEL"

// 	// try logger with a fixed value string, check the
// 	// returned value and the value from the context
// 	lvl = "error"
// 	log = logger(ctx)
// 	leveler = log.Level(&lvl)
// 	if leveler.Level().String() != "ERROR" {
// 		t.Errorf("incorrect level returned from the request, should be ERROR")
// 	}
// 	// pull directly from context
// 	lvl = log.Value(key).(string)
// 	if lvl != "ERROR" {
// 		t.Errorf("incorrect level returned from the context, should be ERROR")
// 	}

// 	// now try via env vars and fetching from context
// 	lvl = "warn"
// 	os.Setenv(key, lvl)
// 	log = logger(ctx)
// 	lvl = log.Value(key).(string)
// 	if lvl != "WARN" {
// 		t.Errorf("incorrect level returned from the context, should be WAN")
// 	}

// }
