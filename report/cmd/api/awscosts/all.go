package awscosts

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awscost"
)

func RegisterGetAwsCostsTop20(log *slog.Logger, conf *config.Config, api huma.API) {
	var operation = huma.Operation{
		OperationID:   "get-awscosts-top20",
		Method:        http.MethodGet,
		Path:          "/v1/awscosts/top20",
		Summary:       "Return top 20 most expensice costs",
		Description:   "Returns a list of most expensive costs",
		DefaultStatus: http.StatusOK,
	}
	huma.Register(api, operation, func(ctx context.Context, input *struct{}) (*GetAwsCostsTop20Response, error) {
		return handleGetAwsCostsTop20(ctx, log, conf, input)
	})
}

// handleGetAwsCostsTop20
func handleGetAwsCostsTop20(ctx context.Context, log *slog.Logger, conf *config.Config, input *struct{}) (response *GetAwsCostsTop20Response, err error) {
	var (
		service *awscost.Service[*AwsCost]
		costs   []*AwsCost
	)
	response = &GetAwsCostsTop20Response{}

	service, err = Service[*AwsCost](ctx, log, conf)
	if err != nil {
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

	return
}
