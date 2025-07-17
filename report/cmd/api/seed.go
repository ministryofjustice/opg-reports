package main

import (
	"context"
	"log/slog"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/seed"
)

// seedDB is called if the database doesnt exist on init, so creates a dummy one
func SeedDB(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	var sqlStore sqlr.RepositoryWriter = sqlr.Default(ctx, log, conf)
	var seedService *seed.Service = seed.Default(ctx, log, conf)
	_, err = seedService.All(sqlStore)
	return
}
