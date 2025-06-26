package awscosts

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awscost"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
)

// Service is a small helper that fetches the service using default data store
func Service[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *awscost.Service[T], err error) {

	srv, err = awscost.Default[T](ctx, log, conf)
	return
}
