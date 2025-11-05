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

// UptimeTableHeaders used for the headings on the tables
type UptimeTableHeaders struct {
	Columns []string `json:"columns"` // the row headings - used in the table body and table header & footer at the start
	Data    []string `json:"data"`    // the data headings - should be the date ranges in order
	Extras  []string `json:"extras"`  // extra headers at the end of the table (trend & total)
}

// UptimeTable contains the uptime data, but formatted as a table for easier front end
// handling (less work per request)
type UptimeTable struct {
	Headers UptimeTableHeaders  `json:"headers"`
	Rows    []map[string]string `json:"rows"`
	Footer  map[string]string   `json:"footer"`
}

// GetAwsUptimeGroupedResponseBody is the response body
type GetAwsUptimeGroupedResponseBody[T api.Model] struct {
	Count   int                    `json:"count"`
	Request *AwsUptimeGroupedInput `json:"request"`
	Dates   []string               `json:"dates"`
	Groups  []string               `json:"groups"`
	Data    []T                    `json:"data"`
	Tabular *UptimeTable           `json:"tabular"` // the data records covnerted to a table structure
}

// GetAwsUptimeGroupedResponse
type GetAwsUptimeGroupedResponse[T api.Model] struct {
	Body GetAwsUptimeGroupedResponseBody[T]
}

// AwsUptimeGroupedInput is the input object for fetching grouped aws costs.
//
// The `Team`  properties are used as a filter. When they are set to _true_
// they are used in the select, group and order by areas of the sql statement.
// When they have any other value they are used as a filter.
//
// This allows the handler to process a range of queries via the same endpoint.
type AwsUptimeGroupedInput struct {
	Granularity string `json:"granularity,omitempty" path:"granularity" default:"monthly" enum:"yearly,monthly" doc:"Determine if the data is grouped by year or month."`
	StartDate   string `json:"start_date,omitempty" path:"start_date" required:"true" doc:"Earliest date to return data from (uses >=). YYYY-MM-DD." example:"2024-03-01" pattern:"([0-9]{4}-[0-9]{2}[\\-0-9]{0,2})"`
	EndDate     string `json:"end_date,omitempty" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}[\\-0-9]{0,2})"`
	Team        string `json:"team,omitempty" query:"team" doc:"Group and filter flag for team. When _true_ adds team.Name to group and selection, when any other value it becomes an exact match filter." example:"true|TeamName"`
	Tabular     bool   `json:"tabular,omitempty" query:"tabular" doc:"When true, tabular version of the data is also included in the response"`
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

	response.Body = GetAwsUptimeGroupedResponseBody[T]{
		Count:   len(data),
		Request: input,
		Dates:   utils.Months(input.StartDate, input.EndDate),
		Groups:  api.GetGroupedByColumns(options),
		Data:    data,
	}

	// If tabular is enabled, convert and configure
	if input.Tabular {
		// handle if no columns are setup
		groups := []string{"date"}
		if len(response.Body.Groups) > 0 {
			groups = response.Body.Groups
		}
		// always add the date column in
		bdy, foot, e := TabulateGroupedUptime(log, groups, response.Body.Dates, data)
		if e == nil {
			response.Body.Tabular = &UptimeTable{
				Headers: UptimeTableHeaders{
					Columns: response.Body.Groups,
					Data:    response.Body.Dates,
					Extras:  []string{"trend", "total"},
				},
				Rows:   bdy,
				Footer: foot,
			}
		}
	}

	return
}
