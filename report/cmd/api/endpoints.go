package main

import (
	"context"
	accountapi "opg-reports/report/internal/domains/account/api"
	pingapi "opg-reports/report/internal/domains/ping/api"
	teamapi "opg-reports/report/internal/domains/team/api"
	"opg-reports/report/packages/httpx"
)

const versionPrefix string = `/v1`

func registerEndpoints(ctx context.Context, mux httpx.MuxServer, cfg httpx.MuxConfigurer) {
	// Home & ping..
	// - home
	httpx.Register(ctx, mux, cfg, `/{$}`, nil, pingapi.Ping)
	// - root ping
	httpx.Register(ctx, mux, cfg, `/ping/{$}`, nil, pingapi.Ping)
	// - versioned home
	httpx.Register(ctx, mux, cfg, versionPrefix+`/{$}`, nil, pingapi.Ping)
	// - versioned ping
	httpx.Register(ctx, mux, cfg, versionPrefix+`/ping/{$}`, nil, pingapi.Ping)

	// Teams
	// - teams list used for navigation
	httpx.Register(ctx, mux, cfg, versionPrefix+`/teams/{$}`, nil, teamapi.TeamsForNavigation)

	// Accounts
	// - list all accounts
	httpx.Register(ctx, mux, cfg, versionPrefix+`/accounts/{$}`, nil, accountapi.Accounts)
	// - list all accounts with team filter
	httpx.Register(ctx, mux, cfg, versionPrefix+`/accounts/team/{team}/{$}`, nil, accountapi.Accounts)

}
