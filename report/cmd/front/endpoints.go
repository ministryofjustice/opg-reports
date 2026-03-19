package main

import (
	"context"
	"opg-reports/report/internal/config"
	static "opg-reports/report/internal/domains/static/home"
	"opg-reports/report/packages/httpx"
)

func registerEndpoints(ctx context.Context, mux httpx.Mux, cfg *config.Config) {

	// deal wtih custom statics
	static.Assets(ctx, mux, cfg)
	static.GovUK(ctx, mux, cfg)
	static.LocalAssets(ctx, mux, cfg)

}
