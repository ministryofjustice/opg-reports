package teams

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// RegisterAllTeams registers the `get-teams-all`
func RegisterAllTeams(log *slog.Logger, conf *config.Config, api huma.API) {
	var operation = huma.Operation{
		OperationID:   "get-teams-all",
		Method:        http.MethodGet,
		Path:          "/v1/teams/all",
		Summary:       "Return all teams",
		Description:   "Returns a list of all teams known about.",
		DefaultStatus: http.StatusOK,
	}
	huma.Register(api, operation, func(ctx context.Context, input *struct{}) (*GetTeamsAllResponse, error) {
		return handleAllTeams(ctx, log, conf, input)
	})
}

// handleAllTeams deals with each request and fetches
func handleAllTeams(ctx context.Context, log *slog.Logger, conf *config.Config, input *struct{}) (response *GetTeamsAllResponse, err error) {
	var (
		service       *team.Service[*team.Team]
		teams         []*team.Team
		responseTeams []*Team = []*Team{}
	)
	response = &GetTeamsAllResponse{}

	service, err = Service(ctx, log, conf)
	if err != nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}

	teams, err = service.GetAllTeams()
	if err != nil {
		err = huma.Error500InternalServerError("failed find all teams", err)
		return
	}
	// convert database team data to response version of Team
	err = utils.Convert(&teams, &responseTeams)
	if err != nil {
		err = huma.Error500InternalServerError("failed to convert data to teams", err)
		return
	}

	response.Body.Data = responseTeams

	return
}
