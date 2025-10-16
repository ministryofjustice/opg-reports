package main

import (
	"log/slog"

	"opg-reports/report/config"

	"github.com/danielgtaylor/huma/v2"
)

// addMiddleware add all middleware into the request;
func addMiddleware(hapi huma.API, log *slog.Logger, conf *config.Config) {

	hapi.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		next(ctx)
	})

}
