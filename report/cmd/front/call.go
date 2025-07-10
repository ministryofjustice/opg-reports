package main

import (
	"context"
	"log/slog"
	"opg-reports/report/config"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/restr"
	"opg-reports/report/internal/service/api"
	"opg-reports/report/internal/service/front"
)

type ApiGetAllTeamsResponse struct {
	Count int         `json:"count,omityempty"`
	Data  []*api.Team `json:"data"`
}

func GetAPITeams(ctx context.Context, log *slog.Logger, conf *config.Config) (teams []*api.Team, err error) {
	var (
		srv    = front.Default[*ApiGetAllTeamsResponse](ctx, log, conf)
		client = restr.Default(ctx, log, conf)
	)
	result, err := srv.GetFromAPI(client, endpoints.TEAMS_GET_ALL)
	if err != nil {
		return
	}
	teams = result.Data
	return
}
