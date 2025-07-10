package teams

import (
	"context"
	"log/slog"
	"net/http"

	"opg-reports/report/config"
	"opg-reports/report/endpoints"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"

	"github.com/danielgtaylor/huma/v2"
)

// GetTeamsAllResponse is response object used by the handler
type GetTeamsAllResponse[T api.Model] struct {
	Body struct {
		Count int `json:"count,omityempty"`
		Data  []T `json:"data"`
	}
}

// RegisterAllTeams registers the `get-teams-all` endpoint
func RegisterGetTeamsAll[T api.Model](log *slog.Logger, conf *config.Config, humaapi huma.API, service api.TeamGetter[T], store sqlr.Reader) {
	var operation = huma.Operation{
		OperationID: "get-teams-all",
		Method:      http.MethodGet,
		Path:        endpoints.TEAMS_GET_ALL,
		// Path:          "/v1/teams/all",
		Summary:       "Return all teams",
		Description:   "Returns a list of all teams known about.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Teams"},
	}
	huma.Register(humaapi, operation, func(ctx context.Context, input *struct{}) (*GetTeamsAllResponse[T], error) {
		return handleGetTeamsAll[T](ctx, log, conf, service, store, input)
	})
}

// handleAllTeams deals with each request and fetches
func handleGetTeamsAll[T api.Model](ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.TeamGetter[T], store sqlr.Reader, input *struct{}) (response *GetTeamsAllResponse[T], err error) {
	var teams []T
	response = &GetTeamsAllResponse[T]{}

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	teams, err = service.GetAllTeams(store)
	if err != nil {
		err = huma.Error500InternalServerError("failed find all teams", err)
		return
	}
	response.Body.Data = teams
	response.Body.Count = len(teams)

	return
}
