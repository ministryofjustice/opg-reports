// Package infracostdiff is very similar to infracostdynamic, but is focused on
// only two months requests (rather than range) and displays the with a change
// higher than the one asked
package infracostdiff

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/utils/ex"
	"opg-reports/report/internal/utils/marshal"
	"opg-reports/report/internal/utils/qb"
	"opg-reports/report/internal/utils/tabulate"
	"opg-reports/report/internal/utils/tabulate/headers"
	"opg-reports/report/internal/utils/tabulate/rows"
	"opg-reports/report/internal/utils/timers"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/infacosts/diff/{date_a}/{date_b}`
	opID          string = `infracosts-get-diff`
	opSummary     string = `Infracost differences table data`
	opDescription string = `Returns a table of infracost data with cost differences between period.`
)

// InfracostDiffRequest is the incoming request options
type InfracostDiffRequest struct {
	DateA   string `path:"date_a" json:"date_a" required:"true" doc:"Start date to use." example:"2026-01" pattern:"([0-9]{4}-[0-9]{2})"` // required - date range input
	DateB   string `path:"date_b" json:"date_b" required:"true" doc:"End date to use." example:"2026-02" pattern:"([0-9]{4}-[0-9]{2})"`   // required - date range input
	Team    string `query:"team" json:"team,omitempty"`
	Account string `query:"account" json:"account,omitempty"`
	Service string `query:"service" json:"service,omitempty"`
	Change  string `query:"change" json:"change" default:"100"`
}

// Months returns the months
func (self *InfracostDiffRequest) Months() (months []string) {
	months = []string{self.DateA, self.DateB}
	return
}

// InfracostDiffResponse is the handlers data struct passed to a huma api which will then be rendered
type InfracostDiffResponse struct {
	Body *InfracostDiffResponseBody
}

// InfracostDiffResponseBody is the response body, containing all data to be returned
type InfracostDiffResponseBody struct {
	Request     *InfracostDiffRequest    `json:"request"`     // the original request
	Headers     map[string][]string      `json:"headers"`     // headers contains details for table headers / rendering
	Data        []map[string]interface{} `json:"data"`        // the actual data results
	Performance []*timers.Timer          `json:"performance"` // duration of the call
	Count       int                      `json:"count"`       // counter to check data aligns
}

// Filter contains all the possible filters passed from the request that arent "true"
type Filter struct {
	Months  []string `db:"months" json:"months"`
	Team    string   `db:"team" json:"team"`
	Account string   `db:"account" json:"account"`
	Service string   `db:"service" json:"service,omitempty"`
}

// querySegments is the possible options to use when query the database
//
// The key should map to the json name in `InfracostDiffRequest`, any `:x`
// values should match the json name in `filter` struct.
//
// Aliases and selected fields should match the json values for the
// returned struct
var querySegments = map[string][]*qb.Segment{
	// for diff, we always want to group by service, account etc for high level
	// of details
	"_default": {
		{Type: qb.SELECT, Stmt: `infracosts.date as date`},
		{Type: qb.SELECT, Stmt: `CAST(COALESCE(SUM(cost), 0) as REAL) as cost`},
		{Type: qb.SELECT, Stmt: `accounts.team_name as team`},
		{Type: qb.SELECT, Stmt: `accounts.id as account_id`},
		{Type: qb.SELECT, Stmt: `accounts.name as account`},
		{Type: qb.SELECT, Stmt: `accounts.environment as environment`},
		{Type: qb.SELECT, Stmt: `infracosts.service as service`},

		{Type: qb.JOIN, Stmt: `LEFT JOIN accounts ON accounts.id = infracosts.account_id`},
		{Type: qb.WHERE, Stmt: `infracosts.service != 'Tax'`},
		{Type: qb.WHERE, Stmt: `infracosts.date IN (:months)`},

		{Type: qb.GROUPBY, Stmt: `infracosts.date`},
		{Type: qb.GROUPBY, Stmt: `infracosts.service`},
		{Type: qb.GROUPBY, Stmt: `accounts.team_name`},
		{Type: qb.GROUPBY, Stmt: `accounts.name`},
		{Type: qb.GROUPBY, Stmt: `accounts.environment`},
	},
	"team": {
		{Type: qb.WHERE, Stmt: `accounts.team_name = :team`},
	},
	"account": {
		{Type: qb.WHERE, Stmt: `accounts.name = :account`},
	},
	"service": {
		{Type: qb.WHERE, Stmt: `infracosts.service = :service`},
	},
}

// the query builder
var builder = qb.New("infracosts", querySegments)

// base table options, add more details to this within the handler
var tableOpts = &tabulate.Options{
	ColumnKey:     "date",
	ValueKey:      "cost",
	RowEndF:       rows.DiffF,
	TableEndF:     tabulate.TotalF,
	TableFilterF:  diffFilterF, // filter the table over 100
	SortByColumn:  "difference",
	SortDirection: "asc",
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

var diffOver float64 = 100

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *InfracostDiffRequest) (*InfracostDiffResponse, error) {
		return getData(timers.ContextWithTimers(ctx), log, db, &operation, input)
	})
}

func getData(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *InfracostDiffRequest) (resp *InfracostDiffResponse, err error) {
	var (
		body        *InfracostDiffResponseBody
		query       *dbstmts.Select[*Filter, *infracostmodels.CostData]
		forFilter   map[string]string
		filter      *Filter                  = &Filter{}
		stmt        string                   = ""
		tableData   []map[string]interface{} = []map[string]interface{}{}
		months      []string                 = []string{}
		requestData map[string]string        = map[string]string{}
		lg          *slog.Logger             = log.With("func", "infracostdiff.getData", "operation", operation.OperationID)
		headings    *headers.Headers         = &headers.Headers{
			Headers: []*headers.Header{
				{Field: "team", Type: headers.KEY, Default: ""},
				{Field: "account_id", Type: headers.KEY, Default: ""},
				{Field: "account", Type: headers.KEY, Default: ""},
				{Field: "environment", Type: headers.KEY, Default: ""},
				{Field: "service", Type: headers.KEY, Default: ""},
				{Field: "difference", Type: headers.END, Default: 0.0},
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

	// generate query statement
	stmt, _ = builder.FromRequest(requestData)
	lg.With("stmt", fmt.Sprintln(stmt)).Debug("sql statement ... ")

	// get months
	months = input.Months()
	// setup headings
	lg.With("headings", headings).Debug("setup headings ...")
	//  - add data headers of the months
	headings.AddDataHeader(months...)

	lg.Debug("creating select statement ...")
	// remove true values from the data for the filter usage
	forFilter = ex.FilterValue(requestData, "true")
	err = marshal.Convert(forFilter, &filter)
	if err != nil {
		return
	}
	// set the min diff level
	if input.Change != "" {
		diffOver, _ = strconv.ParseFloat(input.Change, 64)
	}
	// add the months
	filter.Months = months
	// configure the db query with the generated statement and
	// filter values
	query = &dbstmts.Select[*Filter, *infracostmodels.CostData]{
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

	tableData = tabulate.Tabulate[float64](tableData, headings, tableOpts)

	// prep result
	timers.Stop(ctx, operation.OperationID)
	body = &InfracostDiffResponseBody{
		Request:     input,
		Headers:     headings.ByType(),
		Data:        tableData,
		Count:       len(tableData),
		Performance: timers.All(ctx),
	}
	// setup response
	resp = &InfracostDiffResponse{Body: body}
	lg.Info("complete.")
	return
}

func diffFilterF(table []map[string]interface{}, headings *headers.Headers) []map[string]interface{} {
	var endCol = headings.End()
	var filtered = []map[string]interface{}{}

	for _, row := range table {
		var value = math.Abs(headers.Value[float64](endCol, row))
		if value >= diffOver {
			filtered = append(filtered, row)
		}
	}

	return filtered
}
