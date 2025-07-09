package api

import (
	"context"
	"log/slog"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/seed"
)

func seedDB(ctx context.Context, log *slog.Logger, conf *config.Config) (teams []*sqlr.BoundStatement, awsaccounts []*sqlr.BoundStatement, awscosts []*sqlr.BoundStatement) {

	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	teams, _ = seeder.Teams(sqc)
	awsaccounts, _ = seeder.AwsAccounts(sqc)
	awscosts, _ = seeder.AwsCosts(sqc)

	return
}
