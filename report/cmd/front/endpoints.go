package main

import (
	"context"
	"opg-reports/report/internal/config"
	home "opg-reports/report/internal/domains/home/front"
	static "opg-reports/report/internal/domains/static/home"
	"opg-reports/report/packages/httpx"
)

func registerEndpoints(ctx context.Context, mux httpx.MuxServer, cfg *config.Config) {
	var hook = GetTeamsHook
	// deal wtih custom static registration
	static.Assets(ctx, mux, cfg)
	static.GovUK(ctx, mux, cfg)
	static.LocalAssets(ctx, mux, cfg)
	static.IgnoreFavicon(ctx, mux, cfg)

	// HOME
	// - home
	httpx.Register(ctx, mux, cfg, `/{$}`, hook, home.SetPageBaseline)
}
