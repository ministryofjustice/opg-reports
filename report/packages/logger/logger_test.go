package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
)

func TestPackagesLoggerLevel(t *testing.T) {
	var lvl = "debug"
	var ctx = context.WithValue(t.Context(), ctxLevelKey, "warn")

	// if we set an env var, that should return before the ctx value
	os.Setenv(ctxLevelKey, "error")
	_, l := Level(ctx, &lvl)
	if actual := l.Level().String(); actual != "ERROR" {
		t.Errorf("expected log level to be 'ERROR', actual [%s]", actual)
	}

	// if there is no env var, then it should be a warning
	os.Setenv(ctxLevelKey, "")
	_, l = Level(ctx, &lvl)
	if actual := l.Level().String(); actual != "WARN" {
		t.Errorf("expected log level to be 'warn', actual [%s]", actual)
	}

	// on an empty context with a value passed it should be that value
	ctx = context.TODO()
	_, l = Level(ctx, &lvl)
	if actual := l.Level().String(); actual != "DEBUG" {
		t.Errorf("expected log level to be 'debug', actual [%s]", actual)
	}

	// with a nil, the default should be info
	ctx = context.TODO()
	_, l = Level(ctx, nil)
	if actual := l.Level().String(); actual != "INFO" {
		t.Errorf("expected log level to be 'info', actual [%s]", actual)
	}

}

func TestPackagesLoggerHandler(t *testing.T) {
	var ht = "text"
	var ctx = context.WithValue(t.Context(), ctxHandlerKey, "json")

	// if we set an env var, that should overwrite the json context version
	os.Setenv(ctxHandlerKey, "text")
	_, h := Handler(ctx, &ht, &slog.HandlerOptions{Level: slog.LevelInfo})
	if actual := fmt.Sprintf("%T", h); actual != "*slog.TextHandler" {
		t.Errorf("expected handler to be text, actual [%s]", actual)
	}

	// without the env we should get json from the context
	os.Setenv(ctxHandlerKey, "")
	_, h = Handler(ctx, &ht, &slog.HandlerOptions{Level: slog.LevelInfo})
	if actual := fmt.Sprintf("%T", h); actual != "*slog.JSONHandler" {
		t.Errorf("expected handler to be json, actual [%s]", actual)
	}

	// with empty context should use the parameter
	ctx = context.TODO()
	ht = "json"
	_, h = Handler(ctx, &ht, &slog.HandlerOptions{Level: slog.LevelInfo})
	if actual := fmt.Sprintf("%T", h); actual != "*slog.JSONHandler" {
		t.Errorf("expected handler to be json, actual [%s]", actual)
	}
	// with empty context and nil, should use the default
	ctx = context.TODO()
	_, h = Handler(ctx, nil, &slog.HandlerOptions{Level: slog.LevelInfo})
	if actual := fmt.Sprintf("%T", h); actual != "*slog.TextHandler" {
		t.Errorf("expected handler to be text, actual [%s]", actual)
	}

}
