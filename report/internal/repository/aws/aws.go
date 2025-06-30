package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
)

type Repository struct {
	ctx  context.Context
	conf *config.Config
	log  *slog.Logger
}

func (self *Repository) Close() {
}

func New(ctx context.Context, log *slog.Logger, conf *config.Config) (rp *Repository, err error) {
	rp = &Repository{}

	if log == nil {
		err = fmt.Errorf("no logger passed aws repository")
		return
	}
	if conf == nil {
		err = fmt.Errorf("no config passed aws repository")
		return
	}

	log = log.WithGroup("aws-repository")
	rp = &Repository{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}

	return
}
