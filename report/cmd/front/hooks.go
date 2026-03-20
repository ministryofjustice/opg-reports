package main

import (
	"context"
	"fmt"
	teamapi "opg-reports/report/internal/domains/team/api"
	"opg-reports/report/packages/clients/getter"
	"opg-reports/report/packages/httpx"
	"time"
)

// GetTeamsHook is called by the front end endpoints
// at the end of each request to get consistent data
// instead of replicating the code.
func GetTeamsHook(ctx context.Context, cfg httpx.MuxConfigurer, r httpx.FitleredRequest, resp httpx.MuxResponseType) {
	fmt.Println("called hook!")
	teams, err := GetTeams(ctx, cfg, r)
	if err == nil {
		resp.SetTeams(teams)
	}
}

// GetTeams is a shared helper to getch team data as this is done some all parts of the
// html front end as its the navigation
func GetTeams(ctx context.Context, cfg httpx.MuxConfigurer, r httpx.FitleredRequest) (teams []string, err error) {
	var (
		result = &teamapi.Result{}
		src    = &getter.Source{
			Host:    cfg.ApiHostname(),
			Path:    `/v1/teams/`,
			Timeout: (2 * time.Second),
		}
	)
	// fetch raw data
	result, err = getter.Get[*teamapi.Result](ctx, src, r.RequestData())
	if err != nil {
		return
	}
	// set the team data
	teams = result.Teams
	return
}
