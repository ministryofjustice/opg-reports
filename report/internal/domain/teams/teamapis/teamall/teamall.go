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
type Request struct{}

// Response is the handlers data struct passed to a huma api which will then be rendered
type Response struct {
	Body *ResponseBody
}

// ResponseBody is the response body, containing all data to be returned
type ResponseBody struct {
	Data    []map[string]string `json:"data"`
	Request *Request            `json:"request"`
	Count   int                 `json:"count,omitempty"`
}

// empty is used as the input data for the select statement
// as there are now filters etc
type empty struct{}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {

	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *Request) (*Response, error) {
		return getAllTeams(ctx, log, db, &operation, input)
	})

}

func getAllTeams(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *Request) (response *Response, err error) {
	var (
		body     *ResponseBody
		selector *dbstmts.Select[*empty, *teammodels.Team]
		result   []map[string]string = []map[string]string{}
		lg       *slog.Logger        = log.With("func", "teams.teamapis.teamall.getAllTeams", "operation", operation.OperationID)
	)
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
	body = &ResponseBody{
		Request: input,
		Count:   len(result),
		Data:    result,
	}
	response = &Response{Body: body}
	lg.Info("complete.")
	return
}
