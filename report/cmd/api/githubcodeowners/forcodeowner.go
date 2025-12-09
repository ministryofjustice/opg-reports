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

type GetGithubCodeOwnersForCodeOwnerResponseBody[T api.Model] struct {
	Count   int `json:"count,omityempty"`
	Data    []T `json:"data"`
	Request *GithubCodeOwnersForCodeOwnerInput
}

// GetGithubCodeOwnersForCodeOwnerResponse is response object used by the handler
type GetGithubCodeOwnersForCodeOwnerResponse[T api.Model] struct {
	Body *GetGithubCodeOwnersForCodeOwnerResponseBody[T]
}

type GithubCodeOwnersForCodeOwnerInput struct {
	CodeOwner string `json:"codeowner,omitempty" path:"codeowner" doc:"Filter by this codeowner." example:"ministryofjustice/sirius"`
}

// RegisterGetGithubCodeOwnersForCodeOwner registers the `get-githubcodeowners-for-codeowner` endpoint
func RegisterGetGithubCodeOwnersForCodeOwner[T api.Model](
	log *slog.Logger,
	conf *config.Config,
	humaapi huma.API,
	service api.GithubCodeOwnersForCodeOwnerGetter[T],
	store sqlr.RepositoryReader,
) {
	var operation = huma.Operation{
		OperationID:   "get-githubcodeowners-for-codeowner",
		Method:        http.MethodGet,
		Path:          endpoints.GITHUBCODEOWNERS_FOR_CODEOWNER,
		Summary:       "Return all github codeowner data for this codeowner",
		Description:   "Returns a list of all codeowner data for the codeowner.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Github Code Owners"},
	}
	huma.Register(humaapi, operation, func(ctx context.Context, input *GithubCodeOwnersForCodeOwnerInput) (*GetGithubCodeOwnersForCodeOwnerResponse[T], error) {
		return handleGetGithubCodeOwnersForCodeOwner[T](ctx, log, conf, service, store, input)
	})
}

// handleAllTeams deals with each request and fetches
func handleGetGithubCodeOwnersForCodeOwner[T api.Model](
	ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.GithubCodeOwnersForCodeOwnerGetter[T], store sqlr.RepositoryReader,
	input *GithubCodeOwnersForCodeOwnerInput,
) (response *GetGithubCodeOwnersForCodeOwnerResponse[T], err error) {
	var (
		data []T
		body *GetGithubCodeOwnersForCodeOwnerResponseBody[T]
		opts = &api.GetAllGithubCodeOwnersForCodeOwnerOptions{CodeOwner: strings.ToLower(input.CodeOwner)}
	)
	log.Info("handling get-githubcodeowners-for-codeowner")
	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	data, err = service.GetAllGithubCodeOwnersForCodeOwner(store, opts)
	if err != nil {
		err = huma.Error500InternalServerError("failed to find all github code owners info based on this codeowner", err)
		return
	}
	body = &GetGithubCodeOwnersForCodeOwnerResponseBody[T]{
		Request: input,
		Count:   len(data),
		Data:    data,
	}
	response = &GetGithubCodeOwnersForCodeOwnerResponse[T]{
		Body: body,
	}

	return
}
