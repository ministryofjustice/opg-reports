package ctxlog

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/packages/types"
	"os"
	"testing"
)

var (
	_ types.ContextLogger = &ctxLogger{}
)

func TestPackagesCtxLogCtxLoggerFromSetting(t *testing.T) {

	// setup ctx
	ctx := New(t.Context(), Log(slog.LevelWarn, JSON))
	// check the level
	level := ctx.Logger().Leveler()
	if level != slog.LevelWarn {
		t.Errorf("different log level found that set via env")
	}
	// check handler
	handler := ctx.Logger().Handler()
	if fmt.Sprintf("%T", handler) != "*slog.JSONHandler" {
		t.Errorf("mismatch handler: %T", handler)
	}
	// set a value on context and get it back ...
	newCtx := context.WithValue(ctx, "test-key", "A")
	// now try and get that value...
	val := newCtx.Value("test-key")
	if val != "A" {
		t.Errorf("failed to fetch value correctly ..")
		t.FailNow()
	}
	// now try and get the logger info again
	if ctx.Logger().Leveler() != slog.LevelWarn {
		t.Errorf("log level mismatch")
	}
	if fmt.Sprintf("%T", ctx.Logger().Handler()) != "*slog.JSONHandler" {
		t.Errorf("mismatch handler: %T", handler)
	}

}

func TestPackagesCtxLogCtxLoggerFromEnv(t *testing.T) {

	// test log level and handler from env variables
	os.Setenv("LOG_LEVEL", "ERROR")
	os.Setenv("LOG_TYPE", "JSON")
	// setup ctx
	ctx := New(t.Context(), nil)
	// check the level
	level := ctx.Logger().Leveler()
	if level != slog.LevelError {
		t.Errorf("different log level found that set via env")
	}
	// check handler
	handler := ctx.Logger().Handler()
	if fmt.Sprintf("%T", handler) != "*slog.JSONHandler" {
		t.Errorf("mismatch handler: %T", handler)
	}
	// set a value on context and get it back ...
	newCtx := context.WithValue(ctx, "test-key", "A")
	// now try and get that value...
	val := newCtx.Value("test-key")
	if val != "A" {
		t.Errorf("failed to fetch value correctly ..")
		t.FailNow()
	}
	// now try and get the logger info again
	if ctx.Logger().Leveler() != slog.LevelError {
		t.Errorf("log level mismatch")
	}
	if fmt.Sprintf("%T", ctx.Logger().Handler()) != "*slog.JSONHandler" {
		t.Errorf("mismatch handler: %T", handler)
	}

}
