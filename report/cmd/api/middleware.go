package main

import (
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/config"
)

// addMiddleware add all middleware into the reqquest; currently these are:
//
//   - Check max age of the local database, if older than 3 days, fetch from s3
func addMiddleware(hapi huma.API, log *slog.Logger, conf *config.Config) {
	// add database age information
	hapi.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		downloadLatestDB(ctx.Context(), log, conf)
		next(ctx)
	})

}
