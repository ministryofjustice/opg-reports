package awsuptime

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

// GetAwsUptimeGroupedResponse
type GetAwsUptimeGroupedResponse[T api.Model] struct {
	Body struct {
		Count   int                    `json:"count"`
		Request *AwsUptimeGroupedInput `json:"request"`
		Dates   []string               `json:"dates"`
		Groups  []string               `json:"groups"`
		Data    []T                    `json:"data"`
	}
}

// AwsUptimeGroupedInput is the input object for fetching grouped aws costs.
//
// The `Team`, `Region`, `Service`, `Account` & `Environment` properties are
// used as a filter. When they are set to _true_ they are used in the select,
// group and order by areas of the sql statement. When they have any other
// value they are used as a filter.
//
// This allows the handler to process a large range of AWS cost queries
// via the same endpoint.
type AwsUptimeGroupedInput struct {
	Granularity string `json:"granularity,omitempty" path:"granularity" default:"monthly" enum:"yearly,monthly" doc:"Determine if the data is grouped by year or month."`
	StartDate   string `json:"start_date,omitempty" path:"start_date" required:"true" doc:"Earliest date to return data from (uses >=). YYYY-MM-DD." example:"2024-03-01" pattern:"([0-9]{4}-[0-9]{2}[\\-0-9]{0,2})"`
	EndDate     string `json:"end_date,omitempty" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}[\\-0-9]{0,2})"`
	Team        string `json:"team,omitempty" query:"team" doc:"Group and filter flag for team. When _true_ adds team.Name to group and selection, when any other value it becomes an exact match filter." example:"true|TeamName"`
}

// RegisterGetAwsUptimeGrouped handles all AWS uptime request that are grouped by date + other fields.
func RegisterGetAwsUptimeGrouped[T api.Model](
	log *slog.Logger,
	conf *config.Config,
	humaapi huma.API,
	service api.AwsUptimeGroupedGetter[T],
	store sqlr.RepositoryReader,
) {
	var operation = huma.Operation{
		OperationID:   "get-awsuptime-grouped",
		Method:        http.MethodGet,
		Path:          endpoints.AWSUPTIME_GROUPED,
		Summary:       "Grouped AWS uptime data.",
		Description:   "Returns AWS uptime data grouped by time period and other options.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"AWS Uptime"},
	}
	huma.Register(humaapi, operation, func(ctx context.Context, input *AwsUptimeGroupedInput) (*GetAwsUptimeGroupedResponse[T], error) {
		return handleGetAwsUptimeGrouped(ctx, log, conf, service, store, input)
	})
}

// handleGetAwsUptimeGrouped uses the input parameters to work out a series of options for the
// underlying service to run against the database about how the uptime data is grouped together
// and what fields are used within the select / where etc.
func handleGetAwsUptimeGrouped[T api.Model](
	ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.AwsUptimeGroupedGetter[T], store sqlr.RepositoryReader,
	input *AwsUptimeGroupedInput,
) (response *GetAwsUptimeGroupedResponse[T], err error) {
	var data []T = []T{}
	response = &GetAwsUptimeGroupedResponse[T]{}
	log.Info("handling get-awscosts-grouped")

	if service == nil {
		err = huma.Error500InternalServerError("failed to connect to service", err)
		return
	}
	defer service.Close()

	// if using monthly, set the start & end dates to first and last of month
	if input.Granularity == string(utils.GranularityMonth) {
		input.StartDate, input.EndDate = utils.FirstAndLastOfMonth(input.StartDate, input.EndDate)
	}

	// create the options
	options := &api.GetAwsUptimeGroupedOptions{
		DateFormat: utils.GRANULARITY_TO_FORMAT[input.Granularity],
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		Team:       input.Team,
	}

	data, err = service.GetGroupedAwsUptime(store, options)
	if err != nil {
		err = huma.Error500InternalServerError("failed find grouped costs", err)
		return
	}
	// set the response data
	response.Body.Request = input
	response.Body.Dates = utils.Months(input.StartDate, input.EndDate)
	response.Body.Groups = api.GetGroupedByColumns(options)
	response.Body.Data = data
	response.Body.Count = len(data)

	return
}
