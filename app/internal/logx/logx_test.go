package logx

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
)

// TestLogxGet tries fetching log from context to make sure its found or set
func TestLogxGet(t *testing.T) {
	var (
		exists bool
		ctx    context.Context = context.Background()
	)

	// the logger should be empty
	if ctx.Value(ctxKey) != nil {
		t.Error("logger was alreay set in context")
		t.FailNow()
	}

	// look for existing - should be false
	ctx, _, exists = get(ctx)
	if exists {
		t.Error("logger was found in context before being set")
		t.FailNow()
	}
	// second call should find it as the context should be updated
	ctx, _, exists = get(ctx)
	if !exists {
		t.Errorf("expected to find that current context has a logger")
	}

	// setup context to have a logger
	ctx, _ = New(ctx, nil, nil)
	// now it should be found
	ctx, _, exists = get(ctx)
	if !exists {
		t.Errorf("expected to find logger in the context")
	}

}

func TestLogxLevel(t *testing.T) {
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

func TestLogxHandler(t *testing.T) {
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
