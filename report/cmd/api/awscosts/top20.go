package awscosts

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

// GetAwsCostsTop20Response
type GetAwsCostsTop20Response[T api.Model] struct {
	Body struct {
		Count int `json:"count,omityempty"`
		Data  []T `json:"data"`
	}
}

// RegisterGetAwsCostsTop20 registers the API endpoint to handle returning the 20 most expensive AWS
// costs from the data.
//
// Registers the helper function to handle the request as well as the details for api spec
func RegisterGetAwsCostsTop20[T api.Model](log *slog.Logger, conf *config.Config, humaapi huma.API, service api.AwsCostsTop20Getter[T], store sqlr.RepositoryReader) {
	var operation = huma.Operation{
		OperationID:   "get-awscosts-top20",
		Method:        http.MethodGet,
		Path:          endpoints.AWSCOSTS_GET_TOP_20,
		Summary:       "Top 20 most expensive AWS costs",
		Description:   "Returns a list of most expensive AWS costs stored in the database (excluding tax).",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Costs"},
	}
	huma.Register(humaapi, operation, func(ctx context.Context, input *struct{}) (*GetAwsCostsTop20Response[T], error) {
		return handleGetAwsCostsTop20(ctx, log, conf, service, store, input)
	})
}

// handleGetAwsCostsTop20 returns the 20 most expensive cost entries within the data store.
//
// Tax is excluded, as tax entry is a month's worth of tax in a single entry so will always
// be one of highest line items costs
func handleGetAwsCostsTop20[T api.Model](ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.AwsCostsTop20Getter[T], store sqlr.RepositoryReader, input *struct{}) (response *GetAwsCostsTop20Response[T], err error) {
	var costs []T
	response = &GetAwsCostsTop20Response[T]{}

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	costs, err = service.GetTop20AwsCosts(store)
	if err != nil {
		err = huma.Error500InternalServerError("failed find top 20", err)
		return
	}

	response.Body.Data = costs
	response.Body.Count = len(costs)

	return
}
