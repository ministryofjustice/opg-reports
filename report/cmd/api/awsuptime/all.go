package awsuptime

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

// GetAwsUptimeAllResponse
type GetAwsUptimeAllResponse[T api.Model] struct {
	Body struct {
		Count int `json:"count,omitempty"`
		Data  []T `json:"data"`
	}
}

// RegisterGetAwsUptimeAll handles registering the endpoint to return all aws accounts stored
// in the data.
//
// Attaches both the handler (`handleGetAwsUptimeAll`) and the spec details to the huma api
func RegisterGetAwsUptimeAll[T api.Model](
	log *slog.Logger,
	conf *config.Config,
	humaapi huma.API,
	service api.AwsUptimeGetter[T],
	store sqlr.RepositoryReader,
) {
	var operation = huma.Operation{
		OperationID:   "get-awsuptime-all",
		Method:        http.MethodGet,
		Path:          endpoints.AWSUPTIME_ALL,
		Summary:       "All AWS uptime entries",
		Description:   "Returns a list of all AWS uptime entries stored.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Uptime"},
	}
	huma.Register(humaapi, operation, func(ctx context.Context, input *struct{}) (*GetAwsUptimeAllResponse[T], error) {
		return handleGetAwsUptimeAll(ctx, log, conf, service, store, input)
	})
}

// handleGetAwsAccountsAll deals with each request and fetches
func handleGetAwsUptimeAll[T api.Model](
	ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.AwsUptimeGetter[T], store sqlr.RepositoryReader,
	input *struct{},
) (response *GetAwsUptimeAllResponse[T], err error) {
	var accounts []T
	response = &GetAwsUptimeAllResponse[T]{}
	log.Info("handling get-awsaccounts-all")

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	accounts, err = service.GetAllAwsUptime(store)
	if err != nil {
		err = huma.Error500InternalServerError("failed find uptime entries", err)
		return
	}

	response.Body.Data = accounts
	response.Body.Count = len(accounts)

	return
}
