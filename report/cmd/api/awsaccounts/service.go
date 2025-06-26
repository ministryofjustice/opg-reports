package awsaccounts

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
)

// Service is a small helper that fetches the service for awsaccount.AwsAccount related calls
// and returns that
func Service(ctx context.Context, log *slog.Logger, conf *config.Config) (srv *awsaccount.Service[*awsaccount.AwsAccount], err error) {

	srv, err = awsaccount.Default[*awsaccount.AwsAccount](ctx, log, conf)
	return
}
