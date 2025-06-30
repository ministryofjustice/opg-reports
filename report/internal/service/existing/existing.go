package existing

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
)

const label string = "existing-service"

type Service struct {
	ctx  context.Context
	log  *slog.Logger
	conf *config.Config
}

func New(ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service, err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for %s", label)
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for %s", label)
	}

	srv = &Service{
		ctx:  ctx,
		log:  log.With("service", label),
		conf: conf,
	}
	return
}
