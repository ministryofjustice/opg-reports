package awscosts

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awscost"
)

// GetAwsCostsTop20Response
type GetAwsCostsTop20Response struct {
	Body struct {
		Count int        `json:"count,omityempty"`
		Data  []*AwsCost `json:"data"`
	}
}

// RegisterGetAwsCostsTop20 registers the API endpoint to handle returning the 20 most expensive AWS
// costs from the data.
//
// Registers the helper function to handle the request as well as the details for api spec
func RegisterGetAwsCostsTop20(log *slog.Logger, conf *config.Config, api huma.API, service *awscost.Service[*AwsCost]) {
	var operation = huma.Operation{
		OperationID:   "get-awscosts-top20",
		Method:        http.MethodGet,
		Path:          "/v1/awscosts/top20",
		Summary:       "Top 20 most expensive AWS costs",
		Description:   "Returns a list of most expensive AWS costs stored in the database (excluding tax).",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Costs"},
	}
	huma.Register(api, operation, func(ctx context.Context, input *struct{}) (*GetAwsCostsTop20Response, error) {
		return handleGetAwsCostsTop20(ctx, log, conf, service, input)
	})
}

// handleGetAwsCostsTop20 returns the 20 most expensive cost entries within the data store.
//
// Tax is excluded, as tax entry is a month's worth of tax in a single entry so will always
// be one of highest line items costs
func handleGetAwsCostsTop20(ctx context.Context, log *slog.Logger, conf *config.Config, service *awscost.Service[*AwsCost], input *struct{}) (response *GetAwsCostsTop20Response, err error) {
	var costs []*AwsCost
	response = &GetAwsCostsTop20Response{}

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	costs, err = service.GetTop20()
	if err != nil {
		err = huma.Error500InternalServerError("failed find top 20", err)
		return
	}

	response.Body.Data = costs
	response.Body.Count = len(costs)

	return
}
