package main

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
)

// seedDB is called if the database doesnt exist on init, so creates a dummy one
func SeedDB(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	var sqlStore sqlr.Writer = sqlr.Default(ctx, log, conf)
	var seedService *seed.Service = seed.Default(ctx, log, conf)
	_, err = seedService.All(sqlStore)
	return
}
