package api

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
)

func seedDB(ctx context.Context, log *slog.Logger, conf *config.Config) (teams []*sqlr.BoundStatement, awsaccounts []*sqlr.BoundStatement, awscosts []*sqlr.BoundStatement) {

	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	teams, _ = seeder.Teams(sqc)
	awsaccounts, _ = seeder.AwsAccounts(sqc)
	awscosts, _ = seeder.AwsCosts(sqc)

	return
}
