package infracostsbymonthaccountforteam

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/utils/marshal"
	"opg-reports/report/internal/utils/tabulate"
	"opg-reports/report/internal/utils/tabulate/headers"
	"opg-reports/report/internal/utils/tabulate/rows"
	"opg-reports/report/internal/utils/timers"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/infracosts/between/{start_date}/{end_date}/account/team/{team}`
	opID          string = `infracosts-get-by-month-and-account-for-team`
	opSummary     string = `Return costs grouped by the month and account, filtered by team.`
	opDescription string = `Returns a list costs between the start and end dates grouped by team & account, filtered by team and returned as a table.`
)

const baseSelect string = `
SELECT
    strftime("%Y-%m", infracosts.date) as date,
    CAST(COALESCE(SUM(cost), 0) as REAL) as cost,
	accounts.team_name as team,
	accounts.name as account_name,
	accounts.id as account_id
FROM infracosts
LEFT JOIN accounts ON accounts.id = infracosts.account_id
WHERE
	infracosts.service != 'Tax' AND
	accounts.team_name = :team AND
    strftime("%Y-%m", infracosts.date) IN (:months)
GROUP BY
    accounts.team_name,
	accounts.name,
    strftime("%Y-%m", infracosts.date)
;
`

// filter is used with the sql statment to limit the months to return
type filter struct {
	Months []string `json:"months" db:"months"`
	Team   string   `json:"team" db:"team"`
}

// CostByMonthAccountForTeamRequest contains the incoming url paths and query string data for this endpoint
type CostByMonthAccountForTeamRequest struct {
	StartDate string `json:"start_date" required:"true" path:"start_date" doc:"Earliest date to return data from. YYYY-MM." example:"2025-01" pattern:"([0-9]{4}-[0-9]{2})"`
	EndDate   string `json:"end_date" required:"true" path:"end_date" doc:"Latest date to capture the data for. YYYY-MM."  example:"2025-06" pattern:"([0-9]{4}-[0-9]{2})"`
	Team      string `json:"team" required:"true" path:"team" doc:"team name to filte results by." example:"sirius"`
}

// Start converts the string to a time
func (self *CostByMonthAccountForTeamRequest) Start() (t time.Time) {
	t, _ = times.FromString(self.StartDate)
	return
}

// End converts the string to a time
func (self *CostByMonthAccountForTeamRequest) End() (t time.Time) {
	t, _ = times.FromString(self.EndDate)
	return
}

// CostByMonthAccountForTeamResponse is the handlers data struct passed to a huma api which will then be rendered
type CostByMonthAccountForTeamResponse struct {
	Body *CostByMonthAccountForTeamResponseBody
}

// CostByMonthAccountForTeamResponseBody is the response body, containing all data to be returned
type CostByMonthAccountForTeamResponseBody struct {
	Request     *CostByMonthAccountForTeamRequest `json:"request"`     // the original CostByMonthAccountForTeamRequest
	Headers     map[string][]string               `json:"headers"`     // headers contains details for table headers / rendering
	Data        []map[string]interface{}          `json:"data"`        // the actual data results
	Performance []*timers.Timer                   `json:"performance"` // duration of the call
	Count       int                               `json:"count"`       // counter to check data aligns
}

var tableOpts = &tabulate.Options{
	ColumnKey: "date",
	ValueKey:  "cost",
	RowEndF:   rows.TotalF,
	TableEndF: tabulate.TotalF,
}

// operation describes what this endpoint is doing
var operation = huma.Operation{
	Method:        http.MethodGet,
	DefaultStatus: http.StatusOK,
	Path:          ENDPOINT,
	Summary:       opSummary,
	Description:   opDescription,
	OperationID:   opID,
	Tags:          []string{"infracosts"},
}

// errors
var (
	ErrSelectFailed  = errors.New("select query failed with error.")
	ErrConvertFailed = errors.New("dataset type conversion failed.")
)

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *CostByMonthAccountForTeamRequest) (*CostByMonthAccountForTeamResponse, error) {
		return getByMonthAccountForTeam(timers.ContextWithTimers(ctx), log, db, &operation, input)
	})
}

// getByMonthAccountForTeam fetches the data directly and then converts the db rows from the result into a table row styled
// map to return
func getByMonthAccountForTeam(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *CostByMonthAccountForTeamRequest) (resp *CostByMonthAccountForTeamResponse, err error) {
	var (
		body   *CostByMonthAccountForTeamResponseBody
		query  *dbstmts.Select[*filter, *infracostmodels.CostMonthAccountForTeam]
		lg     *slog.Logger = log.With("func", "infracostsbymonthaccountforteam.getByMonthAccountForTeam", "operation", operation.OperationID)
		months []string     = times.AsYMStrings(times.Months(input.Start(), input.End()))

		tableData []map[string]interface{}
		headings  *headers.Headers = &headers.Headers{Headers: []*headers.Header{
			{Field: "account_name", Type: headers.KEY, Default: ""},
			{Field: "trend", Type: headers.EXTRA, Default: ""},
			{Field: "total", Type: headers.END, Default: 0.0},
		}}
	)
	// timers
	timers.Start(ctx, operation.OperationID)
	defer func() { timers.Stop(ctx) }()

	lg.Info("starting handler ...")
	lg.With("months", months, "team", input.Team).Debug("determined range of months and team ...")
	// add months to the table headers
	headings.AddDataHeader(months...)
	// setup table config - add the month column
	tableOpts.SortByColumn = months[len(months)-1]

	// create the statement
	lg.Debug("creating select statement ...")
	query = &dbstmts.Select[*filter, *infracostmodels.CostMonthAccountForTeam]{
		Statement: baseSelect,
		Data:      &filter{Months: months, Team: input.Team},
	}

	// run the select
	lg.Debug("running select call ...")
	err = dbselects.Select(ctx, log, db, query)
	if err != nil {
		lg.Error("select failed with error", "err", err.Error())
		err = errors.Join(ErrSelectFailed, err)
		return
	}
	// convert data to table form
	err = marshal.Convert(query.Returned, &tableData)
	if err != nil {
		return
	}
	tableData = tabulate.Tabulate[float64](tableData, headings, tableOpts)

	// prep result
	timers.Stop(ctx, operation.OperationID)
	body = &CostByMonthAccountForTeamResponseBody{
		Request:     input,
		Headers:     headings.ByType(),
		Data:        tableData,
		Count:       len(tableData),
		Performance: timers.All(ctx),
	}
	// setup response
	resp = &CostByMonthAccountForTeamResponse{Body: body}
	lg.Info("complete.")
	return
}
