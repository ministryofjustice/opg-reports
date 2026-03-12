package types

import (
	"context"
	"log/slog"
)

type Resetable interface {
	Reset()
}

// Logger interface used to expand on slog
type Logger interface {
	Log() *slog.Logger
	Leveler() slog.Leveler
	Handler() slog.Handler
}

// ContextLogger expands on context and merges in a default slog
// with funcs exposed to easily access the logger
//
// This is used heavily by all main functions to capture the current
// context.
type ContextLogger interface {
	context.Context
	Ctx() context.Context
	Logger() Logger
	Log() *slog.Logger
}

type Contexter interface {
	Ctx() ContextLogger
}

// Contextable ensures strust have a way to set and retrieve the
// relevant context and logger details so they can be used within
// their functions
type Contextable interface {
	Contexter
	SetCtx(ctx ContextLogger)
}
