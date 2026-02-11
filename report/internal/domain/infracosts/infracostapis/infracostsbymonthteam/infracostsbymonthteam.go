package infracostsbymonthteam

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
	ENDPOINT      string = `/v1/infracosts/between/{start_date}/{end_date}/team`
	opID          string = `infracosts-get-by-month-and-team`
	opSummary     string = `Return costs grouped by the month and team.`
	opDescription string = `Returns a list costs between the start and end dates grouped by team and formatted as a table.`
)

// baseSelect - fastest way to get the data
//   - tried multiple selects, one for each month, concurrently - slowest at ~300ms
//   - tried multiple selects, one per month, in sequence - ~150ms
//   - tried sub-selects per month - ~185ms
//   - this one getting row for each month - ~35ms
const baseSelect string = `
SELECT
    strftime("%Y-%m", infracosts.date) as date,
    CAST(COALESCE(SUM(cost), 0) as REAL) as cost,
    accounts.team_name as team
FROM infracosts
LEFT JOIN accounts ON accounts.id = infracosts.account_id
WHERE
	infracosts.service != 'Tax' AND
    strftime("%Y-%m", infracosts.date) IN (:months)
GROUP BY
    accounts.team_name,
    strftime("%Y-%m", infracosts.date)
;
`

// filter is used with the sql statment to limit the months to return
type filter struct {
	Months []string `json:"months" db:"months"`
}

// CostByMonthTeamRequest contains the incoming url paths and query string data for this endpoint
type CostByMonthTeamRequest struct {
	StartDate string `json:"start_date,omitempty" path:"start_date" doc:"Earliest date to return data from. YYYY-MM." example:"2025-01" pattern:"([0-9]{4}-[0-9]{2})"`
	EndDate   string `json:"end_date,omitempty" path:"end_date" doc:"Latest date to capture the data for. YYYY-MM."  example:"2025-06" pattern:"([0-9]{4}-[0-9]{2})"`
}

// Start converts the string to a time
func (self *CostByMonthTeamRequest) Start() (t time.Time) {
	t, _ = times.FromString(self.StartDate)
	return
}

// End converts the string to a time
func (self *CostByMonthTeamRequest) End() (t time.Time) {
	t, _ = times.FromString(self.EndDate)
	return
}

// CostByMonthTeamResponse is the handlers data struct passed to a huma api which will then be rendered
type CostByMonthTeamResponse struct {
	Body *CostByMonthTeamResponseBody
}

// CostByMonthTeamResponseBody is the response body, containing all data to be returned
type CostByMonthTeamResponseBody struct {
	Request     *CostByMonthTeamRequest  `json:"request"`     // the original CostByMonthTeamRequest
	Headers     map[string][]string      `json:"headers"`     // headers contains details for table headers / rendering
	Data        []map[string]interface{} `json:"data"`        // the actual data results
	Performance []*timers.Timer          `json:"performance"` // duration of the call
	Count       int                      `json:"count"`       // counter to check data aligns
}

var tableOpts = &tabulate.Options{
	ColumnKey:    "date",
	ValueKey:     "cost",
	SortByColumn: "team",
	RowEndF:      rows.TotalF,
	TableEndF:    tabulate.TotalF,
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
	huma.Register(humaapi, operation, func(ctx context.Context, input *CostByMonthTeamRequest) (*CostByMonthTeamResponse, error) {
		return getByMonthTeam(timers.ContextWithTimers(ctx), log, db, &operation, input)
	})
}

// getByMonthAndTeam fetches the data directly and then converts the db rows from the result into a table row styled
// map to return
func getByMonthTeam(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *CostByMonthTeamRequest) (resp *CostByMonthTeamResponse, err error) {
	var (
		body   *CostByMonthTeamResponseBody
		query  *dbstmts.Select[*filter, *infracostmodels.CostMonthTeam]
		lg     *slog.Logger = log.With("func", "infracostsbymonthteam.getByMonthTeam", "operation", operation.OperationID)
		months []string     = times.AsYMStrings(times.Months(input.Start(), input.End()))

		tableData []map[string]interface{}
		headings  *headers.Headers = &headers.Headers{Headers: []*headers.Header{
			{Field: "team", Type: headers.KEY, Default: ""},
			{Field: "trend", Type: headers.EXTRA, Default: ""},
			{Field: "total", Type: headers.END, Default: 0.0},
		}}
	)
	// timers
	timers.Start(ctx, operation.OperationID)
	defer func() { timers.Stop(ctx) }()

	lg.Info("starting handler ...")
	lg.With("months", months).Debug("determined range of months ...")
	// add months to the table headers
	headings.AddDataHeader(months...)
	// create the statement
	lg.Debug("creating select statement ...")
	query = &dbstmts.Select[*filter, *infracostmodels.CostMonthTeam]{
		Statement: baseSelect,
		Data:      &filter{Months: months},
	}

	// run the select
	lg.Debug("running select call ...")
	err = dbselects.Select(ctx, log, db, query)
	if err != nil {
		lg.Error("select failed with error", "err", err.Error())
		err = errors.Join(ErrSelectFailed, err)
		return
	}
	// convert to a table format
	err = marshal.Convert(query.Returned, &tableData)
	if err != nil {
		return
	}
	tableData = tabulate.Tabulate[string](tableData, headings, tableOpts)
	// prep result
	timers.Stop(ctx, operation.OperationID)
	body = &CostByMonthTeamResponseBody{
		Request:     input,
		Headers:     headings.ByType(),
		Data:        tableData,
		Count:       len(tableData),
		Performance: timers.All(ctx),
	}
	// setup response
	resp = &CostByMonthTeamResponse{Body: body}
	lg.Info("complete.")
	return
}
