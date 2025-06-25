package awsaccounts

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
)

// Service is a small helper that fetches the service for awsaccount.AwsAccount related calls
// and returns that
func Service(ctx context.Context, log *slog.Logger, conf *config.Config) (srv *awsaccount.Service[*awsaccount.AwsAccount], err error) {
	var datastore *sqldb.Repository[*awsaccount.AwsAccount]

	datastore, err = sqldb.New[*awsaccount.AwsAccount](ctx, log, conf)
	if err != nil {
		return
	}

	srv, err = awsaccount.NewService(ctx, log, conf, datastore)
	return
}
