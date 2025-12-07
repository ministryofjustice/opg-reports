package githubcodeowners

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

type GetGithubCodeOwnersForTeamResponseBody[T api.Model] struct {
	Count   int `json:"count,omityempty"`
	Data    []T `json:"data"`
	Request *GithubCodeOwnersForTeamInput
}

// GetGithubCodeOwnersForTeamResponse is response object used by the handler
type GetGithubCodeOwnersForTeamResponse[T api.Model] struct {
	Body *GetGithubCodeOwnersForTeamResponseBody[T]
}

type GithubCodeOwnersForTeamInput struct {
	Team string `json:"team,omitempty" path:"team" doc:"Filter by this team." example:"TeamName"`
}

// RegisterGetGithubCodeOwnersForTeam registers the `get-githubcodeowners-for-team` endpoint
func RegisterGetGithubCodeOwnersForTeam[T api.Model](
	log *slog.Logger,
	conf *config.Config,
	humaapi huma.API,
	service api.GithubCodeOwnersForTeamGetter[T],
	store sqlr.RepositoryReader,
) {
	var operation = huma.Operation{
		OperationID:   "get-githubcodeowners-for-team",
		Method:        http.MethodGet,
		Path:          endpoints.GITHUBCODEOWNERS_FOR_TEAM,
		Summary:       "Return all github codeowner data for this team",
		Description:   "Returns a list of all codeowner data for the team.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Github Code Owners"},
	}
	huma.Register(humaapi, operation, func(ctx context.Context, input *GithubCodeOwnersForTeamInput) (*GetGithubCodeOwnersForTeamResponse[T], error) {
		return handleGetGithubCodeOwnersForTeam[T](ctx, log, conf, service, store, input)
	})
}

// handleAllTeams deals with each request and fetches
func handleGetGithubCodeOwnersForTeam[T api.Model](
	ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.GithubCodeOwnersForTeamGetter[T], store sqlr.RepositoryReader,
	input *GithubCodeOwnersForTeamInput,
) (response *GetGithubCodeOwnersForTeamResponse[T], err error) {
	var (
		data []T
		body *GetGithubCodeOwnersForTeamResponseBody[T]
		opts = &api.GetAllGithubCodeOwnersForTeamOptions{Team: strings.ToLower(input.Team)}
	)
	log.Info("handling get-githubcodeowners-for-team")
	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	data, err = service.GetAllGithubCodeOwnersForTeam(store, opts)
	if err != nil {
		err = huma.Error500InternalServerError("failed to find all github code owners for team", err)
		return
	}
	body = &GetGithubCodeOwnersForTeamResponseBody[T]{
		Request: input,
		Count:   len(data),
		Data:    data,
	}
	response = &GetGithubCodeOwnersForTeamResponse[T]{
		Body: body,
	}

	return
}
