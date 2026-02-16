package uptimedynamic

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/uptime/uptimemodels"
	"opg-reports/report/internal/utils/ex"
	"opg-reports/report/internal/utils/marshal"
	"opg-reports/report/internal/utils/qb"
	"opg-reports/report/internal/utils/tabulate"
	"opg-reports/report/internal/utils/tabulate/headers"
	"opg-reports/report/internal/utils/tabulate/rows"
	"opg-reports/report/internal/utils/timers"
	"opg-reports/report/internal/utils/times"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/uptime/between/{start_date}/{end_date}`
	opID          string = `uptime-get-dynamic`
	opSummary     string = `Uptime table data`
	opDescription string = `Returns uptime data as a table grouped by the month and filtered based on query.`
)

// UptimeRequest is the incoming request options
type UptimeRequest struct {
	StartDate string `path:"start_date" json:"start_date" required:"true" doc:"Start date to use." example:"2026-01" pattern:"([0-9]{4}-[0-9]{2})"` // required - date range input
	EndDate   string `path:"end_date" json:"end_date" required:"true" doc:"End date to use." example:"2026-02" pattern:"([0-9]{4}-[0-9]{2})"`       // required - date range input
	Team      string `query:"team" json:"team,omitempty"`
	Sort      string `query:"sort" enum:"average,team" json:"-"` // sort data, dont json encode otherwise break the cast to filter
}

// Months returns all months between date range string
func (self *UptimeRequest) Months() (months []string) {
	months = times.AsYMStrings(
		times.Months(times.MustFromString(self.StartDate), times.MustFromString(self.EndDate)))
	return
}

// UptimeResponse is the handlers data struct passed to a huma api which will then be rendered
type UptimeResponse struct {
	Body *UptimeResponseBody
}

// UptimeResponseBody is the response body, containing all data to be returned
type UptimeResponseBody struct {
	Request     *UptimeRequest           `json:"request"`     // the original request
	Headers     map[string][]string      `json:"headers"`     // headers contains details for table headers / rendering
	Data        []map[string]interface{} `json:"data"`        // the actual data results
	Performance []*timers.Timer          `json:"performance"` // duration of the call
	Count       int                      `json:"count"`       // counter to check data aligns
}

// Filter contains all the possible filters passed from the request that arent "true"
type Filter struct {
	Months []string `db:"months" json:"months"`
	Team   string   `db:"team" json:"team"`
}

// querySegments is the possible options to use when query the database
//
// The key should map to the json name in `UptimeRequest`, any `:x`
// values should match the json name in `filter` struct.
//
// Aliases and selected fields should match the json values for the
// returned struct
var querySegments = map[string][]*qb.Segment{
	"_default": {
		{Type: qb.SELECT, Stmt: `uptime.date as date`},
		{Type: qb.SELECT, Stmt: `CAST(COALESCE(AVG(uptime.average), 0) as REAL) as average`},
		{Type: qb.JOIN, Stmt: `LEFT JOIN accounts ON accounts.id = uptime.account_id`},
		{Type: qb.WHERE, Stmt: `uptime.date IN (:months)`},
		{Type: qb.GROUPBY, Stmt: `uptime.date`},
		{Type: qb.ORDERBY, Stmt: `uptime.date ASC`},
	},
	"team": {
		{Type: qb.SELECT, Stmt: `accounts.team_name as team`},
		{Type: qb.WHERE, Stmt: `accounts.team_name = :team`},
		{Type: qb.GROUPBY, Stmt: `accounts.team_name`},
		{Type: qb.ORDERBY, Stmt: `accounts.team_name ASC`},
	},
}

// the query builder
var builder = qb.New("uptime", querySegments)

// base table options, add more details to this within the handler
var tableOpts = &tabulate.Options{
	ColumnKey: "date",
	ValueKey:  "average",
	RowEndF:   rows.AverageF,
	TableEndF: tabulate.AverageF,
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

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *UptimeRequest) (*UptimeResponse, error) {
		return getData(timers.ContextWithTimers(ctx), log, db, &operation, input)
	})
}

func getData(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *UptimeRequest) (resp *UptimeResponse, err error) {
	var (
		body        *UptimeResponseBody
		query       *dbstmts.Select[*Filter, *uptimemodels.UptimeData]
		forFilter   map[string]string
		filter      *Filter                  = &Filter{}
		stmt        string                   = ""
		tableData   []map[string]interface{} = []map[string]interface{}{}
		months      []string                 = []string{}
		requestData map[string]string        = map[string]string{}
		lg          *slog.Logger             = log.With("func", "uptimedynamic.getData", "operation", operation.OperationID)
		headings    *headers.Headers         = &headers.Headers{ // baseline headings, will get expanded from the Request data within the handler
			Headers: []*headers.Header{
				{Field: "trend", Type: headers.EXTRA, Default: ""},
				{Field: "average", Type: headers.END, Default: 0.0},
			},
		}
	)
	// timers
	timers.Start(ctx, operation.OperationID)
	defer func() { timers.Stop(ctx) }()

	lg.With("input", input).Info("starting handler ...")
	// convert input
	err = marshal.Convert(input, &requestData)
	if err != nil {
		return
	}
	// if there is only dates in the request, force teams
	if len(ex.FilterKeys(requestData, "start_date", "end_date")) == 0 {
		input.Team = "true"
		requestData["team"] = "true"
	}
	// generate query statement
	stmt, _ = builder.FromRequest(requestData)
	lg.With("stmt", fmt.Sprintln(stmt)).Debug("sql statement ... ")
	// get months
	months = input.Months()
	// setup headings
	lg.With("headings", headings).Debug("setup headings ...")
	// 	- add headers, exclude the date fields as that becomes the data columns
	headings.AddKeyHeader(requestData, "start_date", "end_date")
	//  - add data headers of the months
	headings.AddDataHeader(months...)

	lg.Debug("creating select statement ...")
	// remove true values from the data for the filter usage
	forFilter = ex.FilterValue(requestData, "true")
	err = marshal.Convert(forFilter, &filter)
	if err != nil {
		return
	}
	// add the months
	filter.Months = months
	// configure the db query with the generated statement and
	// filter values
	query = &dbstmts.Select[*Filter, *uptimemodels.UptimeData]{
		Statement: stmt,
		Data:      filter,
	}

	lg.Debug("running select call ...")
	err = dbselects.Select(ctx, log, db, query)
	if err != nil {
		return
	}
	// convert data to table form
	err = marshal.Convert(query.Returned, &tableData)
	if err != nil {
		return
	}
	// if asked to sort by cost, then actually use the last month
	// otherwise its a string based sort
	if input.Sort == "" {
		input.Sort = "team"
	}
	tableOpts.SortByColumn = input.Sort
	if input.Sort == "average" {
		tableOpts.SortDirection = "desc"
		tableOpts.SortByColumn = months[len(months)-1]
		tableData = tabulate.Tabulate[float64](tableData, headings, tableOpts)
	} else {
		tableOpts.SortDirection = "asc"
		tableData = tabulate.Tabulate[string](tableData, headings, tableOpts)
	}

	// prep result
	timers.Stop(ctx, operation.OperationID)

	body = &UptimeResponseBody{
		Request:     input,
		Headers:     headings.ByType(),
		Data:        tableData,
		Count:       len(tableData),
		Performance: timers.All(ctx),
	}
	// setup response
	resp = &UptimeResponse{Body: body}

	lg.Info("complete.")
	return
}
