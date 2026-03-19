package main

import (
	"context"
	accountapi "opg-reports/report/internal/domains/account/api"
	pingapi "opg-reports/report/internal/domains/ping/api"
	teamapi "opg-reports/report/internal/domains/team/api"
	"opg-reports/report/packages/httpx"
)

const versionPrefix string = `/v1`

func registerEndpoints(ctx context.Context, mux httpx.Mux) {
	// Home & ping..
	// - home
	mux.Register(`/{$}`, pingapi.Ping)
	// - root ping
	mux.Register(`/ping/{$}`, pingapi.Ping)
	// - versioned home
	mux.Register(versionPrefix+`/{$}`, pingapi.Ping)
	// - versioned ping
	mux.Register(versionPrefix+`/ping/{$}`, pingapi.Ping)

	// Teams
	// - teams list used for navigation
	mux.Register(versionPrefix+`/teams/{$}`, teamapi.TeamsForNavigation)

	// Accounts
	// - list all accounts
	mux.Register(versionPrefix+`/accounts/{$}`, accountapi.Accounts)
	// - list all accounts with team filter
	mux.Register(versionPrefix+`/accounts/team/{team}/{$}`, accountapi.Accounts)

}
