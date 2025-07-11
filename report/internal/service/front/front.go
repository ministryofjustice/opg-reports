package front

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/config"
)

const label string = "front-service"

type Response interface{}
type Result interface{}

type Closer interface {
	Close() (err error)
}

type Service[T Response, R Result] struct {
	ctx  context.Context
	log  *slog.Logger
	conf *config.Config
}

// Close function to do any clean up
func (self *Service[T, R]) Close() (err error) {
	return
}

// New tries to create a version of this service using the context, logger and config values
// passed along.
//
// If logger / config is not passed an error is returned
func New[T Response, R Result](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service[T, R], err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for %s", label)
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for %s", label)
	}

	srv = &Service[T, R]{
		ctx:  ctx,
		log:  log.With("service", label),
		conf: conf,
	}
	return
}

// Default creates a service by calling `New` and swallowing any errors.
//
// Errors are logged, but not shared
func Default[T Response, R Result](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service[T, R]) {
	srv, err := New[T, R](ctx, log, conf)
	if err != nil {
		log.Error("error with default", "err", err.Error())
	}
	return
}
