package main

import (
	"context"
	"net/http"
	accountsapi "opg-reports/report/internal/domains/account/api"
	homeapi "opg-reports/report/internal/domains/home/api"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/handler"
)

const version string = `/v1`

// home api endpoints
const (
	home        string = `/{$}`
	homeVersion string = version + home
	ping        string = `/ping/{$}`
	pingVersion string = version + ping
)

// account api endpoints
const (
	accounts     string = version + `/accounts/{$}`
	accountsTeam string = version + `/accounts/team/{team}/{$}`
)

// team api endpoints
const (
	teams string = version + `/teams/{$}`
)

// code base related
const (
	codebases     string = version + `/code/all/{$}`
	codebasesTeam string = version + `/code/all/team/{team}/{$}`
)

// cost related
const (
	costs     string = version + `/costs/all/between/{date_start}/{date_end}/{$}`
	costsTeam string = version + `/costs/all/between/{date_start}/{date_end}/team/{team}{$}`
)

func registerEndpoints(ctx context.Context, mux *http.ServeMux, in *args.API) {

	// HOME / PING
	// - ping and root pages of the api
	handler.RegisterAPI(
		ctx, mux,
		homeapi.Config(in),
		homeapi.Response(in),
		homeapi.T(),
		home, homeVersion, ping, pingVersion,
	)
	// ACCOUNTS
	// - all accounts, or accounts filtered by team name
	handler.RegisterAPI(
		ctx, mux,
		accountsapi.Config(in),
		accountsapi.Response(in),
		accountsapi.T(),
		accounts, accountsTeam,
	)

	// handler.RegisterAPI(ctx, mux, homeapi.Handler(in), home, homeVersion, ping, pingVersion)
	// handler.RegisterAPI(ctx, mux, accountapi.Handler(in), accounts, accountsTeam)
	// // TEAMS
	// // - all teams
	// handler.RegisterAPI(ctx, mux, teamapi.Handler(in), teams)
	// // CODEBASES
	// // - list all codebases; currently no team filter
	// handler.RegisterAPI(ctx, mux, codebasesapi.Handler(in), codebases)

	// COSTS
	// - costs grouped by month and team
	// handler.RegisterAPI(ctx, mux, costsbyteamapi.Handler(in), costs, codebasesTeam)
}
