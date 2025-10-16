package api

import (
	"context"
	"log/slog"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/seed"
)

func seedDB(ctx context.Context, log *slog.Logger, conf *config.Config) (results *seed.SeedAllResults, err error) {
	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	results, err = seeder.All(sqc)
	return
}
