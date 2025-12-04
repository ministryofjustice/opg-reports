package githubcodeowners

import (
	"context"
	"log/slog"
	"net/http"

	"opg-reports/report/config"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"

	"github.com/danielgtaylor/huma/v2"
)

// GetGithubCodeOwnersAllResponse is response object used by the handler
type GetGithubCodeOwnersAllResponse[T api.Model] struct {
	Body struct {
		Count int `json:"count,omityempty"`
		Data  []T `json:"data"`
	}
}

// RegisterGetGithubCodeOwnersAll registers the `get-githubcodeowners-all` endpoint
func RegisterGetGithubCodeOwnersAll[T api.Model](
	log *slog.Logger,
	conf *config.Config,
	humaapi huma.API,
	service api.GithubCodeOwnersGetter[T],
	store sqlr.RepositoryReader,
) {
	var operation = huma.Operation{
		OperationID:   "get-githubcodeowners-all",
		Method:        http.MethodGet,
		Path:          endpoints.TEAMS_GET_ALL,
		Summary:       "Return all github codeowners",
		Description:   "Returns a list of all codeowners known about.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"GithubCodeOwners"},
	}
	huma.Register(humaapi, operation, func(ctx context.Context, input *struct{}) (*GetGithubCodeOwnersAllResponse[T], error) {
		return handleGetGithubCodeOwnersAll[T](ctx, log, conf, service, store, input)
	})
}

// handleAllTeams deals with each request and fetches
func handleGetGithubCodeOwnersAll[T api.Model](
	ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.GithubCodeOwnersGetter[T], store sqlr.RepositoryReader,
	input *struct{},
) (response *GetGithubCodeOwnersAllResponse[T], err error) {
	var data []T
	response = &GetGithubCodeOwnersAllResponse[T]{}

	log.Info("handling get-githubcodeowners-all")

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	data, err = service.GetAllGithubCodeOwners(store)
	if err != nil {
		err = huma.Error500InternalServerError("failed find all github code owners", err)
		return
	}
	response.Body.Data = data
	response.Body.Count = len(data)

	return
}
