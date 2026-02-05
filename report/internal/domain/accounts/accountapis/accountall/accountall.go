package accountall

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/utils/marshal"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/accounts/all`
	opID          string = `accounts-get-all`
	opSummary     string = `Return all accounts and their teams.`
	opDescription string = `Returns a list of all accounts and team data from the database without filtering.`
)

// operation describes what this endpoint is doing
var operation = huma.Operation{
	Method:        http.MethodGet,
	DefaultStatus: http.StatusOK,
	Path:          ENDPOINT,
	Summary:       opSummary,
	Description:   opDescription,
	OperationID:   opID,
	Tags:          []string{"accounts"},
}

// errors
var (
	ErrSelectFailed  = errors.New("select query failed with error.")
	ErrConvertFailed = errors.New("dataset type conversion failed.")
)

// selectStmt is the main sql statment to fetch data from the db
var selectStmt string = `
SELECT
	id,
	name,
	label,
	environment,
	team_name
FROM accounts
ORDER BY
	team_name,
	name,
	environment ASC
;`

// Request contains the incoming url and query string data for this endpoint
type Request struct{}

// Response is the handlers data struct passed to a huma api which will then be rendered
type Response struct {
	Body *ResponseBody
}

// ResponseBody is the response body, containing all data to be returned
type ResponseBody struct {
	Request     *Request            `json:"request"`
	Performance *perf               `json:"performance"` // duration of the request
	Data        []map[string]string `json:"data"`
	Count       int                 `json:"count,omitempty"`
}

// pref tracks performance of the request to this endpoint, logging start & endtime
// as well as the duration from starting the handler till finishing
type perf struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration string    `json:"duration"`
}

// empty is used as the input data for the select statement
// as there are now filters etc
type empty struct{}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {

	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *Request) (*Response, error) {
		return getAllAccounts(ctx, log, db, &operation, input)
	})

}

// getAllAccounts
func getAllAccounts(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *Request) (response *Response, err error) {
	var (
		body      *ResponseBody
		selector  *dbstmts.Select[*empty, *accountmodels.AccountRow]
		callEnd   time.Time
		callStart time.Time           = time.Now().UTC()
		result    []map[string]string = []map[string]string{}
		lg        *slog.Logger        = log.With("func", "domain.teamapis.teamall.getAllTeams", "operation", operation.OperationID)
	)
	lg.Info("starting handler ...")
	// create the statement
	lg.Debug("creating select statement ...")
	selector = &dbstmts.Select[*empty, *accountmodels.AccountRow]{
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
	callEnd = time.Now().UTC()
	body = &ResponseBody{
		Request: input,
		Count:   len(result),
		Data:    result,
		Performance: &perf{
			Start:    callStart,
			End:      callEnd,
			Duration: fmt.Sprintf("%v", callEnd.Sub(callStart).String()),
		},
	}
	response = &Response{Body: body}
	lg.Info("complete.")
	return
}
