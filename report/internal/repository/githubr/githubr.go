package githubr

import (
	"context"
	"fmt"
	"log/slog"

	"opg-reports/report/config"
)

const label string = "github-repository"

type Repository struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository, err error) {

	if log == nil {
		err = fmt.Errorf("no logger passed for %s", label)
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed for %s", label)
		return
	}

	rp = &Repository{
		ctx:  ctx,
		log:  log.WithGroup(label),
		conf: conf,
	}

	return
}

func Default(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository) {
	rp, err := New(ctx, log, conf)
	if err != nil {
		log.Error("error with default", "err", err.Error())
	}
	return
}
