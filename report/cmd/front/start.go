package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/service/api"
	"opg-reports/report/internal/service/front"
)

type apiTeams struct {
	Count int         `json:"count,omityempty"`
	Data  []*api.Team `json:"data"`
}

func parseTeamList(response *apiTeams) (teams []string, err error) {
	teams = []string{}
	for _, team := range response.Data {
		if team.Name != "Legacy" && team.Name != "ORG" {
			teams = append(teams, team.Name)
		}
	}
	return
}

// registerHandlersFunc type aliassing for easier def
type registerHandlersFunc func(ctx context.Context, log *slog.Logger, conf *config.Config, info *FrontInfo, mux *http.ServeMux)

func StartServer(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	info *FrontInfo,
	mux *http.ServeMux,
	server *http.Server,
	registerFuncs ...registerHandlersFunc,
) {
	var (
		srv = front.Default[*apiTeams, []string](ctx, log, conf)
		ep  = endpoints.TEAMS_GET_ALL
	)
	// get all teams when starting and attach the names
	// - ignore the error in case the api is down
	Info.Teams, _ = srv.GetFromAPI(info.RestClient, ep, parseTeamList)
	// call each register function
	var addr = server.Addr
	for _, registerF := range registerFuncs {
		registerF(ctx, log, conf, info, mux)
	}

	log.Info("Starting front server...")
	log.Info(fmt.Sprintf("FRONT: [http://%s/]", addr))

	server.ListenAndServe()
}
