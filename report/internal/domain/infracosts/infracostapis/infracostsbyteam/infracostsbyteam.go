package infracostsbyteam

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/utils/times"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/infracosts/by-month-and-team/{start_date}/{end_date}`
	opID          string = `infracosts-get-by-team-and-month`
	opSummary     string = `Return costs grouped by the month and team.`
	opDescription string = `Returns a list costs between the start and end dates grouped by team and formatted as a table.`
)

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

// Request contains the incoming url paths and query string data for this endpoint
type Request struct {
	StartDate string `json:"start_date,omitempty" path:"start_date" doc:"Earliest date to return data from (uses >=). YYYY-MM." example:"2025-01" pattern:"([0-9]{4}-[0-9]{2})"`
	EndDate   string `json:"end_date,omitempty" path:"end_date" doc:"Latest date to capture the data for (uses <). YYYY-MM."  example:"2025-06" pattern:"([0-9]{4}-[0-9]{2})"`
}

// Start converts the string to a time
func (self *Request) Start() (t time.Time) {
	t, _ = times.FromString(self.StartDate)
	return
}

// End converts the string to a time
func (self *Request) End() (t time.Time) {
	t, _ = times.FromString(self.EndDate)
	return
}

// Response is the handlers data struct passed to a huma api which will then be rendered
type Response struct {
	Body *ResponseBody
}

// ResponseBody is the response body, containing all data to be returned
type ResponseBody struct {
	Request     *Request            `json:"request"`     // the original request
	Months      []string            `json:"months"`      // months within the range specified
	Headers     *headers            `json:"headers"`     // headers contains details for table headers / rendering
	Data        []map[string]string `json:"data"`        // the actual data results
	Performance *perf               `json:"performance"` // duration of the request
	Count       int                 `json:"count"`       // counter to check data aligns
}

// track the headings to use in the table for easier rendering
type headers struct {
	Labels  []string `json:"labels"`  // labels at the start of a row (table headers)
	Columns []string `json:"columns"` // the core data of the table row (monthly totals)
	Extras  []string `json:"extras"`  // additional columns and the end of row - like row totals
}

// pref tracks performance of the request to this endpoint, logging start & endtime
// as well as the duration from starting the handler till finishing
type perf struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration string    `json:"duration"`
}

// empty is used as there are no placeholders within the sql,
// its fully generated from the data ranges creating multiple
// sub selects
type empty struct{}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *Request) (*Response, error) {
		return getInfracostsGroupedByMonthAndTeam(ctx, log, db, &operation, input)
	})

}

// getInfracostsGroupedByMonthAndTeam
func getInfracostsGroupedByMonthAndTeam(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *Request) (response *Response, err error) {
	var (
		body      *ResponseBody
		callEnd   time.Time
		selector  *dbstmts.Select[*empty, map[string]interface{}]
		callStart time.Time           = time.Now().UTC()
		result    []map[string]string = []map[string]string{}
		months    []string            = []string{}
		lg        *slog.Logger        = log.With("func", "infracostsbyteam.getInfracostsGroupedByMonthAndTeam", "operation", operation.OperationID)
	)
	lg.Info("starting handler ...")

	// work out months in time frame, including the last month
	months = times.AsYMStrings(times.Months(input.Start(), input.End()))
	lg.With("months", months).Debug("determined range of months ...")

	// create the statement
	lg.Debug("creating select statement ...")
	selector = &dbstmts.Select[*empty, map[string]interface{}]{
		Statement: selectStmt(months),
		Data:      &empty{},
	}

	// run the select
	lg.Debug("running select call ...")
	err = dbselects.SelectMap(ctx, log, db, selector)
	if err != nil {
		lg.Error("select failed with error", "err", err.Error())
		err = errors.Join(ErrSelectFailed, err)
		return
	}

	// clean up the returned data to remove _ prefix on dynamic columns
	for _, row := range selector.Returned {
		var entry = map[string]string{"trend": ""}
		for k, v := range row {
			entry[strings.TrimPrefix(k, "_")] = fmt.Sprintf("%v", v)
		}
		result = append(result, entry)
	}

	// prep result
	callEnd = time.Now().UTC()
	body = &ResponseBody{
		Request: input,
		Months:  months,
		Count:   len(result),
		Data:    result,
		Performance: &perf{
			Start:    callStart,
			End:      callEnd,
			Duration: fmt.Sprintf("%v", callEnd.Sub(callStart).String()),
		},
		Headers: &headers{
			Labels:  []string{"team"},
			Columns: months,
			Extras:  []string{"trend", "total"},
		},
	}
	response = &Response{Body: body}
	lg.Info("complete.")
	return
}

// baseSelect is the outline of the select to fetch and join over time period
const baseSelect string = `
SELECT
{subs}
{total}
accA.team_name as team
FROM infracosts as A
LEFT JOIN accounts as accA ON accA.id = A.account_id
GROUP BY
accA.team_name
;`

// selectSub contains the repeating select statment we use for each month to fetch the data
// and generate a table row style output
const selectSub string = `
(
	SELECT
		COALESCE(SUM(cost), 0) as cost
	FROM infracosts
	LEFT JOIN accounts ON accounts.id = infracosts.account_id
	WHERE
		strftime("%Y-%m", infracosts.date) = '{month}' AND
		accounts.team_name = accA.team_name
	GROUP BY
		accounts.team_name, strftime("%Y-%m", infracosts.date)
) as '_{month}',
`

// totalSub works out row rotals for the table so they can be added to the end
const totalSub string = `
(
	SELECT
		COALESCE(SUM(cost), 0) as cost
	FROM infracosts
	LEFT JOIN accounts ON accounts.id = infracosts.account_id
	WHERE
		strftime("%Y-%m", infracosts.date) IN ({monthList}) AND
		accounts.team_name = accA.team_name
	GROUP BY
		accounts.team_name
) as 'total',
`

// selectStmt generates the select from multiple sub-selects so the end result
// has months as column headers, making rendering much easier
func selectStmt(months []string) (stmt string) {
	var (
		baseStmt    string = baseSelect
		subKey      string = `{subs}`
		totalKey    string = `{total}`
		selects     string = ""
		monthString string = fmt.Sprintf("'%s'", strings.Join(months, "','"))
	)
	// iterate over the months to create the joins and sub selects
	for _, month := range months {
		selects += strings.ReplaceAll(selectSub, "{month}", month)
	}
	// now replace the total query with real values
	baseStmt = strings.ReplaceAll(baseStmt, totalKey, strings.ReplaceAll(totalSub, "{monthList}", monthString))
	stmt = strings.ReplaceAll(baseStmt, subKey, selects)

	fmt.Println(stmt)
	return

}
