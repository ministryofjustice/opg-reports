package teamall

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/teams/teammodels"
	"opg-reports/report/internal/utils/marshal"
	"opg-reports/report/internal/utils/timers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/teams/all`
	opID          string = `teams-get-all`
	opSummary     string = `Return all teams.`
	opDescription string = `Returns a list of all teams from the database without filtering.`
)

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

// errors
var (
	ErrSelectFailed  = errors.New("select query failed with error.")
	ErrConvertFailed = errors.New("dataset type conversion failed.")
)

// selectStmt is the main sql statment to fetch data from the db
var selectStmt string = `
SELECT
	name
FROM teams
ORDER BY
	name ASC
;`

// Request contains the incoming url and query string data for this endpoint
type TeamRequest struct{}

// Response is the handlers data struct passed to a huma api which will then be rendered
type TeamResponse struct {
	Body *TeamResponseBody
}

// ResponseBody is the response body, containing all data to be returned
type TeamResponseBody struct {
	Data        []map[string]string `json:"data"`
	Request     *TeamRequest        `json:"request"`
	Performance []*timers.Timer     `json:"performance"`
	Count       int                 `json:"count,omitempty"`
}

// empty is used as the input data for the select statement
// as there are now filters etc
type empty struct{}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {

	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *TeamRequest) (*TeamResponse, error) {
		return getAllTeams(timers.ContextWithTimers(ctx), log, db, &operation, input)
	})

}

// getAllTeams
func getAllTeams(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *TeamRequest) (response *TeamResponse, err error) {
	var (
		body     *TeamResponseBody
		selector *dbstmts.Select[*empty, *teammodels.Team]
		result   []map[string]string = []map[string]string{}
		lg       *slog.Logger        = log.With("func", "teamall.getAllTeams", "operation", operation.OperationID)
	)
	// timers
	timers.Start(ctx, operation.OperationID)
	defer func() { timers.Stop(ctx) }()

	lg.Info("starting handler ...")
	// create the statement
	lg.Debug("creating select statement ...")
	selector = &dbstmts.Select[*empty, *teammodels.Team]{
		Statement: selectStmt,
		Data:      &empty{},
	}
	// run the select
	lg.Debug("running select call ...")
	err = dbselects.Select(ctx, log, db, selector)
	if err != nil {
		lg.Error("select failed with error", "err", err.Error())
		err = errors.Join(ErrSelectFailed, err)
		return
	}
	// handle converting the results to the outbound format by using covnert on the
	// result struct
	err = marshal.Convert(selector.Returned, &result)
	if err != nil {
		lg.Error("marshal type conversion failed", "err", err.Error())
		err = errors.Join(ErrConvertFailed, err)
		return
	}
	// prep result
	timers.Stop(ctx, operation.OperationID)
	body = &TeamResponseBody{
		Request:     input,
		Data:        result,
		Count:       len(result),
		Performance: timers.All(ctx),
	}
	response = &TeamResponse{Body: body}
	lg.Info("complete.")
	return
}
