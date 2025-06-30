package awscosts

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/service/awscost"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// GetAwsGroupedCostsResponse
type GetAwsGroupedCostsResponse struct {
	Body struct {
		Count int               `json:"count,omityempty"`
		Data  []*AwsGroupedCost `json:"data"`
	}
}

// GroupedAwsCostsInput is the input object for fetching grouped aws costs.
//
// The `Team`, `Region`, `Service`, `Account` & `Environment` properties are
// used as a filter. When they are set to <true> they are used in the select,
// group and order by areas of the sql statement. When they have any other
// value they are used as a filter.
//
// This allows the handler to process a large range of AWS cost queries
// via the same endpoint.
type GroupedAwsCostsInput struct {
	Granularity string `json:"granularity,omitempty" path:"granularity" default:"month" enum:"year,month,day" doc:"Determine if the data is grouped by year, month or day."`
	StartDate   string `json:"start_date,omitempty" path:"start_date" required:"true" doc:"Earliest date to return data from (uses >=). YYYY-MM-DD." example:"2024-03-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
	EndDate     string `json:"end_date,omitempty" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}-[0-9]{2})"`

	Team        string `json:"team,omitempty" query:"team" doc:"Group and filter flag for team. When <true> adds team.Name to group and selection, when any other value it becomes an exact match filter."`
	Region      string `json:"region,omitempty" query:"region" doc:"Group and filter flag for AWS region. When <true> adds AWS region to group and selection, when any other value it becomes an exact match filter."`
	Service     string `json:"service,omitempty" query:"service" doc:"Group and filter flag for AWS service name. When <true> adds service name to group and selection, when any other value it becomes an exact match filter."`
	Account     string `json:"account,omitempty" query:"account" doc:"Group and filter flag for AWS account ID. When <true> adds AWS account ID to group and selection, when any other value it becomes an exact match filter."`
	Environment string `json:"environment,omitempty" query:"environment" doc:"Group and filter flag for account environment type. When <true> adds environment to group and selection, when any other value it becomes an exact match filter."`
}

// RegisterGetAwsGroupedCosts handles all AWS Cost request that are grouped by date + other fields.
func RegisterGetAwsGroupedCosts(log *slog.Logger, conf *config.Config, api huma.API, service *awscost.Service[*AwsGroupedCost]) {
	var operation = huma.Operation{
		OperationID:   "get-awscosts-grouped",
		Method:        http.MethodGet,
		Path:          "/v1/awscosts/grouped/{granularity}/{start_date}/{end_date}",
		Summary:       "Grouped AWS cost data.",
		Description:   "Returns AWS costs data (excluding tax) grouped by time period and other options.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Costs"},
	}
	huma.Register(api, operation, func(ctx context.Context, input *GroupedAwsCostsInput) (*GetAwsGroupedCostsResponse, error) {
		return handleGetAwsGroupedCosts(ctx, log, conf, service, input)
	})
}

// handleGetAwsGroupedCosts uses the input parameters to work out a series of options for the
// underlying service to run against the database about how the cost data is grouped together
// and what fields are used within the select / where etc.
func handleGetAwsGroupedCosts(ctx context.Context, log *slog.Logger, conf *config.Config, service *awscost.Service[*AwsGroupedCost], input *GroupedAwsCostsInput) (response *GetAwsGroupedCostsResponse, err error) {
	var (
		costs []*AwsGroupedCost = []*AwsGroupedCost{}
	)
	response = &GetAwsGroupedCostsResponse{}

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	// create the options
	options := &awscost.GetGroupedCostsOptions{
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		DateFormat:  utils.GRANULARITY_TO_FORMAT[input.Granularity],
		Team:        utils.TrueOrFilter(input.Team),
		Region:      utils.TrueOrFilter(input.Region),
		Service:     utils.TrueOrFilter(input.Service),
		Account:     utils.TrueOrFilter(input.Account),
		Environment: utils.TrueOrFilter(input.Environment),
	}

	costs, err = service.GetGroupedCosts(options)
	if err != nil {
		err = huma.Error500InternalServerError("failed find grouped costs", err)
		return
	}

	response.Body.Data = costs
	response.Body.Count = len(costs)

	return
}
