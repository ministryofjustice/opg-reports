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

// CostTableHeaders used for the headings on the tables
type CostTableHeaders struct {
	Columns []string `json:"columns"` // the row headings - used in the table body and table header & footer at the start
	Data    []string `json:"data"`    // the data headings - should be the date ranges in order
	Extras  []string `json:"extras"`  // extra headers at the end of the table (trend & total)
}

// CostTable contains the cost data, but formatted as a table for easier front end handling (less work per request)
type CostTable struct {
	Headers CostTableHeaders    `json:"headers"`
	Rows    []map[string]string `json:"rows"`
	Footer  map[string]string   `json:"footer"`
}

// GetAwsCostsGroupedResponseBody is the body of the response returned
type GetAwsCostsGroupedResponseBody[T api.Model] struct {
	Count   int                   `json:"count"`   // number of db results from the query
	Request *AwsCostsGroupedInput `json:"request"` // the original request for comparison
	Dates   []string              `json:"dates"`   // the date range the request is for
	Groups  []string              `json:"groups"`  // the data grouping that the request generated
	Data    []T                   `json:"data"`    // the data records
	Tabular *CostTable            `json:"tabular"` // the data records covnerted to a table structure
}

// GetAwsCostsGroupedResponse - this is the main struct returned by the handlers
type GetAwsCostsGroupedResponse[T api.Model] struct {
	Body GetAwsCostsGroupedResponseBody[T]
}

// AwsCostsGroupedInput is the input object for fetching grouped aws costs.
//
// The `Team`, `Region`, `Service`, `Account` & `Environment` properties are
// used as a filter. When they are set to _true_ they are used in the select,
// group and order by areas of the sql statement. When they have any other
// value they are used as a filter.
//
// This allows the handler to process a large range of AWS cost queries
// via the same endpoint.
type AwsCostsGroupedInput struct {
	Granularity string `json:"granularity,omitempty" path:"granularity" default:"monthly" enum:"yearly,monthly" doc:"Determine if the data is grouped by year or month."`
	StartDate   string `json:"start_date,omitempty" path:"start_date" required:"true" doc:"Earliest date to return data from (uses >=). YYYY-MM-DD." example:"2024-03-01" pattern:"([0-9]{4}-[0-9]{2}[\\-0-9]{0,2})"`
	EndDate     string `json:"end_date,omitempty" path:"end_date" required:"true" doc:"Latest date to capture the data for (uses <). YYYY-MM-DD."  example:"2024-04-01" pattern:"([0-9]{4}-[0-9]{2}[\\-0-9]{0,2})"`

	Team        string `json:"team,omitempty" query:"team" doc:"Group and filter flag for team. When _true_ adds team.Name to group and selection, when any other value it becomes an exact match filter." example:"true|TeamName"`
	Region      string `json:"region,omitempty" query:"region" doc:"Group and filter flag for AWS region. When _true_ adds AWS region to group and selection, when any other value it becomes an exact match filter." example:"true|eu-west-1"`
	Service     string `json:"service,omitempty" query:"service" doc:"Group and filter flag for AWS service name. When _true_ adds service name to group and selection, when any other value it becomes an exact match filter." example:"true|ECS"`
	Account     string `json:"account,omitempty" query:"account" doc:"Group and filter flag for AWS account ID. When _true_ adds AWS account ID to group and selection, when any other value it becomes an exact match filter." example:"true|1000000001"`
	Label       string `json:"account_label,omitempty" query:"account_label" doc:"Group and filter flag for AWS account label. When _true_ adds AWS account label to group and selection, when any other value it becomes an exact match filter." example:"true|service-name"`
	AccountName string `json:"account_name,omitempty" query:"account_name" doc:"Group and filter flag for AWS account name. When _true_ adds AWS account name to group and selection, when any other value it becomes an exact match filter." example:"true|Account Name"`
	Environment string `json:"environment,omitempty" query:"environment" enum:"development,preproduciton,production,backup,true" doc:"Group and filter flag for account environment type. When _true_ adds environment to group and selection, when any other value it becomes an exact match filter."`

	Tabular bool `json:"tabular,omitempty" query:"tabular" doc:"When true, tabular version of the data is also included in the response"`
}

// RegisterGetAwsCostsGrouped handles all AWS Cost request that are grouped by date + other fields.
func RegisterGetAwsCostsGrouped[T api.Model](
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
	huma.Register(humaapi, operation, func(ctx context.Context, input *AwsCostsGroupedInput) (*GetAwsCostsGroupedResponse[T], error) {
		return handleGetAwsCostsGrouped(ctx, log, conf, service, store, input)
	})
}

// handleGetAwsCostsGrouped uses the input parameters to work out a series of options for the
// underlying service to run against the database about how the cost data is grouped together
// and what fields are used within the select / where etc.
func handleGetAwsCostsGrouped[T api.Model](
	ctx context.Context, log *slog.Logger, conf *config.Config,
	service api.AwsCostsGroupedGetter[T], store sqlr.RepositoryReader,
	input *AwsCostsGroupedInput,
) (response *GetAwsCostsGroupedResponse[T], err error) {
	var costs []T = []T{}

	log = log.With("operation", "handleGetAwsCostsGrouped")
	log.With("input", input).Info("handling get-awscosts-grouped")
	response = &GetAwsCostsGroupedResponse[T]{}

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
	options := &api.GetAwsCostsGroupedOptions{
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
	response.Body = GetAwsCostsGroupedResponseBody[T]{
		Count:   len(costs),
		Request: input,
		Dates:   utils.Months(input.StartDate, input.EndDate),
		Groups:  api.GetGroupedByColumns(options),
		Data:    costs,
	}
	// If tabular is enabled, convert and configure
	if input.Tabular {
		// handle if no columns are setup
		groups := []string{"date"}
		if len(response.Body.Groups) > 0 {
			groups = response.Body.Groups
		}
		// always add the date column in
		bdy, foot, e := TabulateGroupedCosts(log, groups, response.Body.Dates, costs)
		if e == nil {
			response.Body.Tabular = &CostTable{
				Headers: CostTableHeaders{
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
