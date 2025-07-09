package main

import (
	"log/slog"

	"opg-reports/report/config"

	"github.com/danielgtaylor/huma/v2"
)

// addMiddleware add all middleware into the reqquest; currently these are:
//
//   - Check max age of the local database, if older than 3 days, fetch from s3
func addMiddleware(hapi huma.API, log *slog.Logger, conf *config.Config) {
	// check database age
	hapi.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		downloadLatestDB(ctx.Context(), log, conf)
		next(ctx)
	})

}
