package awsaccounts

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
)

// Service is a small helper that fetches the service for awsaccount.AwsAccount related calls
// and returns that
func Service[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *awsaccount.Service[T], err error) {

	srv, err = awsaccount.Default[T](ctx, log, conf)
	return
}
