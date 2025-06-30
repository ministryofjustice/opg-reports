package teams

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/service/team"
)

// RegisterAllTeams registers the `get-teams-all` endpoint
func RegisterGetTeamsAll(log *slog.Logger, conf *config.Config, api huma.API, service *team.Service[*Team]) {
	var operation = huma.Operation{
		OperationID:   "get-teams-all",
		Method:        http.MethodGet,
		Path:          "/v1/teams/all",
		Summary:       "Return all teams",
		Description:   "Returns a list of all teams known about.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Teams"},
	}
	huma.Register(api, operation, func(ctx context.Context, input *struct{}) (*GetTeamsAllResponse, error) {
		return handleGetTeamsAll(ctx, log, conf, service, input)
	})
}

// handleAllTeams deals with each request and fetches
func handleGetTeamsAll(ctx context.Context, log *slog.Logger, conf *config.Config, service *team.Service[*Team], input *struct{}) (response *GetTeamsAllResponse, err error) {
	var (
		teams []*Team
	)
	response = &GetTeamsAllResponse{}

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	teams, err = service.GetAllTeams()
	if err != nil {
		err = huma.Error500InternalServerError("failed find all teams", err)
		return
	}
	response.Body.Data = teams
	response.Body.Count = len(teams)

	return
}
