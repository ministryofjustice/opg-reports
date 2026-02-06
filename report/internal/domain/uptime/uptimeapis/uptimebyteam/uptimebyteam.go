package uptimebyteam

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
	"opg-reports/report/internal/utils/timers"
	"opg-reports/report/internal/utils/times"
	"sort"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/uptime/by-month-and-team/{start_date}/{end_date}`
	opID          string = `uptime-get-by-team-and-month`
	opSummary     string = `Return uptime grouped by team and month.`
	opDescription string = `Returns uptime data grouped by team name and the year-month date.`
)

// baseSelect - fastest way to get the data
const baseSelect string = `
SELECT
    strftime("%Y-%m", uptime.date) as date,
	CAST(COALESCE(AVG(uptime.average), 0) as REAL) as average,
    accounts.team_name
FROM uptime
LEFT JOIN accounts ON accounts.id = uptime.account_id
WHERE
    strftime("%Y-%m", uptime.date) IN (:months)
GROUP BY
    accounts.team_name,
    strftime("%Y-%m", uptime.date)
;
`

// UptimeByMonthTeamRequest contains the incoming url paths and query string data for this endpoint
type UptimeByMonthTeamRequest struct {
	StartDate string `json:"start_date,omitempty" path:"start_date" doc:"Earliest date to return data from (uses >=). YYYY-MM." example:"2025-01" pattern:"([0-9]{4}-[0-9]{2})"`
	EndDate   string `json:"end_date,omitempty" path:"end_date" doc:"Latest date to capture the data for (uses <). YYYY-MM."  example:"2025-06" pattern:"([0-9]{4}-[0-9]{2})"`
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
	Request     *UptimeByMonthTeamRequest `json:"request"`     // the original request
	Months      []string                  `json:"months"`      // months within the range specified
	Headers     *UptimeByMonthTeamHeaders `json:"headers"`     // headers contains details for table headers / rendering
	Data        []map[string]interface{}  `json:"data"`        // the actual data results
	Performance []*timers.Timer           `json:"performance"` // duration of the request
	Count       int                       `json:"count"`       // counter to check data aligns

}

// track the headings to use in the table for easier rendering
type UptimeByMonthTeamHeaders struct {
	Labels  []string `json:"labels"`  // labels at the start of a row (table headers)
	Columns []string `json:"columns"` // the core data of the table row (monthly totals)
	Extras  []string `json:"extras"`  // additional columns between columns & totals
	End     []string `json:"end"`     //
}

// Textual returns headers that should be strings
func (self *UptimeByMonthTeamHeaders) Textual() (list []string) {
	list = []string{}
	list = append(list, self.Labels...)
	list = append(list, self.Extras...)
	return
}
func (self *UptimeByMonthTeamHeaders) Numeric() (list []string) {
	list = []string{}
	list = append(list, self.Columns...)
	list = append(list, self.End...)
	return
}
func (self *UptimeByMonthTeamHeaders) DataColumns() (list []string) {
	return self.Columns
}

// empty is used as there are no placeholders within the sql,
// its fully generated from the data ranges creating multiple
// sub selects
type empty struct{}

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

// used for table rows
var (
	headerLabels  = []string{"team"}
	headerExtras  = []string{"trend"}
	headerOverall = []string{"overall"}
)

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *UptimeByMonthTeamRequest) (*UptimeByMonthTeamResponse, error) {
		return getByMonthTeam(ctx, log, db, &operation, input)
	})

}

// getByMonthTeam
func getByMonthTeam(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *UptimeByMonthTeamRequest) (resp *UptimeByMonthTeamResponse, err error) {
	var (
		body             *UptimeByMonthTeamResponseBody
		table            []map[string]interface{} // table formatted version of the database rows
		query            *dbstmts.Select[*empty, *uptimemodels.UptimeMonthTeam]
		lg               *slog.Logger              = log.With("func", "uptime.getByMonthTeam", "operation", operation.OperationID)
		monthStr, months                           = times.JoinedYMList(times.Months(input.Start(), input.End()))
		headers          *UptimeByMonthTeamHeaders = &UptimeByMonthTeamHeaders{
			Labels:  headerLabels,
			Columns: months,
			Extras:  headerExtras,
			End:     headerOverall,
		}
	)
	timers.Start(operation.OperationID)
	defer func() { timers.Stop(operation.OperationID) }()

	lg.Info("starting handler ...")
	lg.With("months", months).Debug("determined range of months ...")

	// create the statement
	lg.Debug("creating select statement ...")
	query = &dbstmts.Select[*empty, *uptimemodels.UptimeMonthTeam]{
		Statement: strings.ReplaceAll(baseSelect, ":months", monthStr),
		Data:      &empty{},
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
	table = tabular(query.Returned, headers)
	// prep result
	timers.Stop(operation.OperationID)
	body = &UptimeByMonthTeamResponseBody{
		Request:     input,
		Months:      months,
		Count:       len(table),
		Data:        table,
		Performance: timers.AllTimers(),
		Headers:     headers,
	}
	// setup response
	resp = &UptimeByMonthTeamResponse{Body: body}
	lg.Info("complete.")
	return
}

// rowKey used by tabulation to decide the key for each row in the table
func rowKey(row map[string]interface{}) string {
	return strings.ToLower(row["team_name"].(string))
}

// updates the labels on the row to values from the database
func rowLabelFunc(dbRow map[string]interface{}, tableRow map[string]interface{}, headers tabulate.TableHeaders) map[string]interface{} {
	tableRow["team"] = dbRow["team_name"]
	return tableRow
}

// rowUpdatefunc used by tabulation to update each row with average uptime
func rowUpdatefunc(dbRow map[string]interface{}, tableRow map[string]interface{}, headers tabulate.TableHeaders) map[string]interface{} {
	var month = dbRow["date"].(string)
	tableRow[month] = dbRow["average"].(float64)
	return tableRow
}

// rowAverageFunc updates the table rows average values based on count and total sum
func rowAverageFunc(dbRow map[string]interface{}, tableRow map[string]interface{}, headers tabulate.TableHeaders) map[string]interface{} {
	var (
		cols    []string = headers.DataColumns()
		total   float64  = 0.0
		average float64  = 0.0
		count   int      = len(cols)
	)
	for _, col := range cols {
		total += tableRow[col].(float64)
	}
	average = total / float64(count)
	tableRow["overall"] = average

	return tableRow
}

// tableSort - sort the table by the team name for consistency
func tableSort(table []map[string]interface{}, headers tabulate.TableHeaders) []map[string]interface{} {
	sort.Slice(table, func(i, j int) bool {
		var a = table[i]["team"].(string)
		var b = table[j]["team"].(string)
		return (a < b)
	})
	return table
}

// tableSummaary runs over all of the table creating a new row based on the total average of each column
func tableSummaary(table []map[string]interface{}, headers tabulate.TableHeaders) []map[string]interface{} {
	var newRow = tabulate.SkeletonRow(headers)
	var cols = headers.DataColumns()
	var count = len(table)
	var total float64 = 0.0
	// generate total for each column
	for _, row := range table {
		for _, col := range cols {
			newRow[col] = newRow[col].(float64) + row[col].(float64)
		}
		total += row["overall"].(float64)
	}
	// deal with dividing by count for average of averages
	for _, col := range cols {
		newRow[col] = newRow[col].(float64) / float64(count)
	}
	newRow["team"] = "Overall"
	newRow["overall"] = total / float64(count)
	table = append(table, newRow)

	return table
}

// tabular wraps around the main tabular helper to create the return data
func tabular(results []*uptimemodels.UptimeMonthTeam, headers *UptimeByMonthTeamHeaders) (table []map[string]interface{}) {
	var dbRows = []map[string]interface{}{}
	var opts = &tabulate.TabulateOptions{
		Headers:    headers,
		KeyF:       rowKey,
		LabelF:     rowLabelFunc,
		ColumnF:    rowUpdatefunc,
		RowEndF:    rowAverageFunc,
		TableSortF: tableSort,
		TableEndF:  tableSummaary,
	}
	marshal.Convert(results, &dbRows)
	table = tabulate.Tabulate(dbRows, opts)
	return
}
