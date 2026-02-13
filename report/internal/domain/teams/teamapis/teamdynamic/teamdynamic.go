package teamdynamic

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/teams/teammodels"
	"opg-reports/report/internal/utils/ex"
	"opg-reports/report/internal/utils/marshal"
	"opg-reports/report/internal/utils/qb"
	"opg-reports/report/internal/utils/timers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/teams`
	opID          string = `teams-get-dynamic`
	opSummary     string = `Teams listing.`
	opDescription string = `Returns a list of all teams.`
)

// TeamRequest is the incoming request options
type TeamRequest struct {
	Sort string `query:"sort" json:"-"` // sort data, dont json encode otherwise break the cast to filter
}

// TeamResponse is the handlers data struct passed to a huma api which will then be rendered
type TeamResponse struct {
	Body *TeamResponseBody
}

// TeamResponseBody is the response body, containing all data to be returned
type TeamResponseBody struct {
	Request     *TeamRequest       `json:"request"`     // the original request
	Data        []*teammodels.Team `json:"data"`        // the actual data results
	Performance []*timers.Timer    `json:"performance"` // duration of the call
	Count       int                `json:"count"`       // counter to check data aligns
}

// Filter contains all the possible filters passed from the request that arent "true"
// - currently empty
type Filter struct {
}

// querySegments is the possible options to use when query the database
//
// The key should map to the json name in `TeamRequest`, any `:x`
// values should match the json name in `filter` struct.
//
// Aliases and selected fields should match the json values for the
// returned struct
var querySegments = map[string][]*qb.Segment{
	"_default": {
		{Type: qb.SELECT, Stmt: `teams.name`},
	},
}

// the query builder
var builder = qb.New("teams", querySegments)

// operation describes what this endpoint is doing
var operation = huma.Operation{
	Method:        http.MethodGet,
	DefaultStatus: http.StatusOK,
	Path:          ENDPOINT,
	Summary:       opSummary,
	Description:   opDescription,
	OperationID:   opID,
	Tags:          []string{"teams"},
}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *TeamRequest) (*TeamResponse, error) {
		return getData(timers.ContextWithTimers(ctx), log, db, &operation, input)
	})
}

func getData(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *TeamRequest) (resp *TeamResponse, err error) {
	var (
		body        *TeamResponseBody
		query       *dbstmts.Select[*Filter, *teammodels.Team]
		forFilter   map[string]string
		filter      *Filter           = &Filter{}
		stmt        string            = ""
		requestData map[string]string = map[string]string{}
		lg          *slog.Logger      = log.With("func", "teamdynamic.getData", "operation", operation.OperationID)
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
	query = &dbstmts.Select[*Filter, *teammodels.Team]{
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
	body = &TeamResponseBody{
		Request:     input,
		Data:        query.Returned,
		Count:       len(query.Returned),
		Performance: timers.All(ctx),
	}
	// setup response
	resp = &TeamResponse{Body: body}
	lg.Info("complete.")
	return
}
