package uptimebymonthteam

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/uptime/uptimemodels"
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
	ENDPOINT      string = `/v1/uptime/between/{start_date}/{end_date}/team`
	opID          string = `uptime-get-by-month-and-team`
	opSummary     string = `Return uptime grouped by team and month.`
	opDescription string = `Returns uptime data grouped by team name and the year-month date.`
)

// baseSelect - fastest way to get the data
const baseSelect string = `
SELECT
    strftime("%Y-%m", uptime.date) as date,
	CAST(COALESCE(AVG(uptime.average), 0) as REAL) as average,
    accounts.team_name as team
FROM uptime
LEFT JOIN accounts ON accounts.id = uptime.account_id
WHERE
    strftime("%Y-%m", uptime.date) IN (:months)
GROUP BY
    accounts.team_name,
    strftime("%Y-%m", uptime.date)
;
`

// filter is used with the sql statment to limit the months to return
type filter struct {
	Months []string `json:"months" db:"months"`
}

// UptimeByMonthTeamRequest contains the incoming url paths and query string data for this endpoint
type UptimeByMonthTeamRequest struct {
	StartDate string `json:"start_date,omitempty" path:"start_date" doc:"Earliest date to return data from. YYYY-MM." example:"2025-11" pattern:"([0-9]{4}-[0-9]{2})"`
	EndDate   string `json:"end_date,omitempty" path:"end_date" doc:"Latest date to capture the data for. YYYY-MM."  example:"2026-01" pattern:"([0-9]{4}-[0-9]{2})"`
}

// Start converts the string to a time
func (self *UptimeByMonthTeamRequest) Start() (t time.Time) {
	t, _ = times.FromString(self.StartDate)
	return
}

// End converts the string to a time
func (self *UptimeByMonthTeamRequest) End() (t time.Time) {
	t, _ = times.FromString(self.EndDate)
	return
}

// UptimeByMonthTeamResponse is the handlers data struct passed to a huma api which will then be rendered
type UptimeByMonthTeamResponse struct {
	Body *UptimeByMonthTeamResponseBody
}

// UptimeByMonthTeamResponseBody is the response body, containing all data to be returned
type UptimeByMonthTeamResponseBody struct {
	Request     *UptimeByMonthTeamRequest `json:"request"` // the original request
	Headers     map[string][]string       `json:"headers"` // headers contains details for table headers / rendering
	Data        []map[string]interface{}  `json:"data"`    // the actual data results
	Count       int                       `json:"count"`   // counter to check data aligns
	Performance []*timers.Timer           `json:"performance"`
}

var tableOpts = &tabulate.Options{
	ColumnKey:    "date",
	ValueKey:     "average",
	SortByColumn: "team",
	RowEndF:      rows.AverageF,
	TableEndF:    tabulate.AverageF,
}

// operation describes what this endpoint is doing
var operation = huma.Operation{
	Method:        http.MethodGet,
	DefaultStatus: http.StatusOK,
	Path:          ENDPOINT,
	Summary:       opSummary,
	Description:   opDescription,
	OperationID:   opID,
	Tags:          []string{"uptime"},
}

// errors
var (
	ErrSelectFailed  = errors.New("select query failed with error.")
	ErrConvertFailed = errors.New("dataset type conversion failed.")
)

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *UptimeByMonthTeamRequest) (*UptimeByMonthTeamResponse, error) {
		return getByMonthTeam(timers.ContextWithTimers(ctx), log, db, &operation, input)
	})

}

// getByMonthTeam
func getByMonthTeam(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *UptimeByMonthTeamRequest) (resp *UptimeByMonthTeamResponse, err error) {
	var (
		body   *UptimeByMonthTeamResponseBody
		query  *dbstmts.Select[*filter, *uptimemodels.UptimeMonthTeam]
		lg     *slog.Logger = log.With("func", "uptime.getByMonthTeam", "operation", operation.OperationID)
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
	query = &dbstmts.Select[*filter, *uptimemodels.UptimeMonthTeam]{
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
	body = &UptimeByMonthTeamResponseBody{
		Request:     input,
		Headers:     headings.ByType(),
		Data:        tableData,
		Count:       len(tableData),
		Performance: timers.All(ctx),
	}
	// setup response
	resp = &UptimeByMonthTeamResponse{Body: body}
	lg.Info("complete.")
	return
}
