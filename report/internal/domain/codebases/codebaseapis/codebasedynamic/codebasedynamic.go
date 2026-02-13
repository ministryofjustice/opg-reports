package codebasedynamic

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/utils/ex"
	"opg-reports/report/internal/utils/marshal"
	"opg-reports/report/internal/utils/qb"
	"opg-reports/report/internal/utils/timers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/codebases`
	opID          string = `codebases-get-dynamic`
	opSummary     string = `Codebase listing.`
	opDescription string = `Returns a list of all codebases and codeowner data.`
)

// CodebaseRequest is the incoming request options
type CodebaseRequest struct {
	Team string `query:"team" json:"team,omitempty"`
	Sort string `query:"sort" json:"-"` // sort data, dont json encode otherwise break the cast to filter
}

// CodebaseResponse is the handlers data struct passed to a huma api which will then be rendered
type CodebaseResponse struct {
	Body *CodebaseResponseBody
}

// CodebaseResponseBody is the response body, containing all data to be returned
type CodebaseResponseBody struct {
	Request     *CodebaseRequest               `json:"request"`     // the original request
	Data        []*codebasemodels.CodebaseData `json:"data"`        // the actual data results
	Performance []*timers.Timer                `json:"performance"` // duration of the call
	Count       int                            `json:"count"`       // counter to check data aligns
}

// Filter contains all the possible filters passed from the request that arent "true"
type Filter struct {
	Team string `db:"team" json:"team"`
}

// querySegments is the possible options to use when query the database
//
// The key should map to the json name in `CodebaseRequest`, any `:x`
// values should match the json name in `filter` struct.
//
// Aliases and selected fields should match the json values for the
// returned struct
var querySegments = map[string][]*qb.Segment{
	"_default": {
		{Type: qb.SELECT, Stmt: `codebases.name`},
		{Type: qb.SELECT, Stmt: `codebases.full_name`},
		{Type: qb.SELECT, Stmt: `codebases.url`},
		{Type: qb.SELECT, Stmt: `json_group_array( DISTINCT json_object( 'name', codeowners.name, 'team_name', codeowners.team_name ) ) filter ( where codeowners.name is not null) as codeowner_list`},
		{Type: qb.JOIN, Stmt: `LEFT JOIN codeowners on codeowners.codebase_full_name = codebases.full_name`},
		{Type: qb.GROUPBY, Stmt: `codebases.full_name`},
		{Type: qb.ORDERBY, Stmt: `codeowners.team_name ASC`},
		{Type: qb.ORDERBY, Stmt: `codebases.full_name ASC`},
	},
	"team": {
		{Type: qb.WHERE, Stmt: `codeowners.team_name = :team`},
	},
}

// the query builder
var builder = qb.New("codebases", querySegments)

// operation describes what this endpoint is doing
var operation = huma.Operation{
	Method:        http.MethodGet,
	DefaultStatus: http.StatusOK,
	Path:          ENDPOINT,
	Summary:       opSummary,
	Description:   opDescription,
	OperationID:   opID,
	Tags:          []string{"codebases"},
}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *CodebaseRequest) (*CodebaseResponse, error) {
		return getData(timers.ContextWithTimers(ctx), log, db, &operation, input)
	})
}

func getData(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *CodebaseRequest) (resp *CodebaseResponse, err error) {
	var (
		body        *CodebaseResponseBody
		query       *dbstmts.Select[*Filter, *codebasemodels.CodebaseData]
		forFilter   map[string]string
		filter      *Filter           = &Filter{}
		stmt        string            = ""
		requestData map[string]string = map[string]string{}
		lg          *slog.Logger      = log.With("func", "codebasedynamic.getData", "operation", operation.OperationID)
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

	lg.Debug("creating select statement ...")
	// remove true values from the data for the filter usage
	forFilter = ex.FilterValue(requestData, "true")
	err = marshal.Convert(forFilter, &filter)
	if err != nil {
		return
	}
	// configure the db query with the generated statement and
	// filter values
	query = &dbstmts.Select[*Filter, *codebasemodels.CodebaseData]{
		Statement: stmt,
		Data:      filter,
	}
	lg.Debug("running select call ...")
	err = dbselects.Select(ctx, log, db, query)
	if err != nil {
		return
	}

	// prep result
	timers.Stop(ctx, operation.OperationID)
	body = &CodebaseResponseBody{
		Request:     input,
		Data:        query.Returned,
		Count:       len(query.Returned),
		Performance: timers.All(ctx),
	}
	// setup response
	resp = &CodebaseResponse{Body: body}
	lg.Info("complete.")
	return
}
