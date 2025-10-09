package awscosts

import (
	"context"
	"log/slog"
	"net/http"

	"opg-reports/report/config"
	"opg-reports/report/internal/endpoints"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"
	"opg-reports/report/internal/utils"

	"github.com/danielgtaylor/huma/v2"
)

// GetAwsGroupedCostsResponse
type GetAwsGroupedCostsResponse[T api.Model] struct {
	Body struct {
		Count   int                   `json:"count"`
		Request *GroupedAwsCostsInput `json:"request"`
		Dates   []string              `json:"dates"`
		Groups  []string              `json:"groups"`
		Data    []T                   `json:"data"`
	}
}

// GroupedAwsCostsInput is the input object for fetching grouped aws costs.
//
// The `Team`, `Region`, `Service`, `Account` & `Environment` properties are
// used as a filter. When they are set to _true_ they are used in the select,
// group and order by areas of the sql statement. When they have any other
// value they are used as a filter.
//
// This allows the handler to process a large range of AWS cost queries
// via the same endpoint.
type GroupedAwsCostsInput struct {
	Granularity string `json:"granularity,omitempty" path:"granularity" default:"monthly" enum:"yearly,monthly" doc:"Determine if the data is grouped by year, month or day."`
	StartDate   string `json:"start_date,omitempty" path:"start_date" required:"true" doc:"Earliest date to return data from (uses >=). YYYY-MM-DD." example:"2024-03-01" pattern:"([0-9]{4}-[0-9]{2}[\\-0-9]{0,2})"`
	EndDate     string `json:"end_date,omitempty" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}[\\-0-9]{0,2})"`

	Team        string `json:"team,omitempty" query:"team" doc:"Group and filter flag for team. When _true_ adds team.Name to group and selection, when any other value it becomes an exact match filter." example:"true|TeamName"`
	Region      string `json:"region,omitempty" query:"region" doc:"Group and filter flag for AWS region. When _true_ adds AWS region to group and selection, when any other value it becomes an exact match filter." example:"true|eu-west-1"`
	Service     string `json:"service,omitempty" query:"service" doc:"Group and filter flag for AWS service name. When _true_ adds service name to group and selection, when any other value it becomes an exact match filter." example:"true|ECS"`
	Account     string `json:"account,omitempty" query:"account" doc:"Group and filter flag for AWS account ID. When _true_ adds AWS account ID to group and selection, when any other value it becomes an exact match filter." example:"true|1000000001"`
	Label       string `json:"account_label,omitempty" query:"account_label" doc:"Group and filter flag for AWS account label. When _true_ adds AWS account label to group and selection, when any other value it becomes an exact match filter." example:"true|service-name"`
	AccountName string `json:"account_name,omitempty" query:"account_name" doc:"Group and filter flag for AWS account name. When _true_ adds AWS account name to group and selection, when any other value it becomes an exact match filter." example:"true|Account Name"`
	Environment string `json:"environment,omitempty" query:"environment" enum:"development,preproduciton,production,backup,true" doc:"Group and filter flag for account environment type. When _true_ adds environment to group and selection, when any other value it becomes an exact match filter."`
}

// RegisterGetAwsGroupedCosts handles all AWS Cost request that are grouped by date + other fields.
func RegisterGetAwsGroupedCosts[T api.Model](
	log *slog.Logger,
	conf *config.Config,
	humaapi huma.API,
	service api.AwsCostsGroupedGetter[T],
	store sqlr.RepositoryReader,
) {
	var operation = huma.Operation{
		OperationID:   "get-awscosts-grouped",
		Method:        http.MethodGet,
		Path:          endpoints.AWSCOSTS_GROUPED,
		Summary:       "Grouped AWS cost data.",
		Description:   "Returns AWS costs data (excluding tax) grouped by time period and other options.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Costs"},
	}
	huma.Register(humaapi, operation, func(ctx context.Context, input *GroupedAwsCostsInput) (*GetAwsGroupedCostsResponse[T], error) {
		return handleGetAwsGroupedCosts(ctx, log, conf, service, store, input)
	})
}

// handleGetAwsGroupedCosts uses the input parameters to work out a series of options for the
// underlying service to run against the database about how the cost data is grouped together
// and what fields are used within the select / where etc.
func handleGetAwsGroupedCosts[T api.Model](
	ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.AwsCostsGroupedGetter[T], store sqlr.RepositoryReader,
	input *GroupedAwsCostsInput,
) (response *GetAwsGroupedCostsResponse[T], err error) {
	var costs []T = []T{}
	response = &GetAwsGroupedCostsResponse[T]{}
	log.Info("handling get-awscosts-grouped")

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	// if using monthly, set the start & end dates to first and last of month
	if input.Granularity == string(utils.GranularityMonth) {
		if sd, e := utils.StringToTime(input.StartDate); e == nil {
			input.StartDate = utils.FirstDayOfMonth(sd).Format(utils.DATE_FORMATS.YMD)
		}
		if ed, e := utils.StringToTime(input.EndDate); e == nil {
			input.EndDate = utils.LastDayOfMonth(ed).Format(utils.DATE_FORMATS.YMD)
		}
	}

	// create the options
	options := &api.GetGroupedCostsOptions{
		DateFormat:  utils.GRANULARITY_TO_FORMAT[input.Granularity],
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Team:        input.Team,
		Region:      input.Region,
		Service:     input.Service,
		Account:     input.Account,
		AccountName: input.AccountName,
		Environment: input.Environment,
		Label:       input.Label,
	}

	costs, err = service.GetGroupedAwsCosts(store, options)
	if err != nil {
		err = huma.Error500InternalServerError("failed find grouped costs", err)
		return
	}
	// set the response data
	response.Body.Request = input
	response.Body.Dates = utils.Months(input.StartDate, input.EndDate)
	response.Body.Groups = options.Groups()
	response.Body.Data = costs
	response.Body.Count = len(costs)

	return
}
