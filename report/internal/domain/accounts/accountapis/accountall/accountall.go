package accountall

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/utils/timers"

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
type AccountRequest struct{}

// Response is the handlers data struct passed to a huma api which will then be rendered
type AccountResponse struct {
	Body *AccountResponseBody
}

// ResponseBody is the response body, containing all data to be returned
type AccountResponseBody struct {
	Request     *AccountRequest             `json:"request"`
	Performance []*timers.Timer             `json:"performance"` // duration of the request
	Data        []*accountmodels.AccountRow `json:"data"`
	Count       int                         `json:"count,omitempty"`
}

// empty is used as the input data for the select statement
// as there are now filters etc
type empty struct{}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {

	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *AccountRequest) (*AccountResponse, error) {
		return getAllAccounts(ctx, log, db, &operation, input)
	})

}

// getAllAccounts
func getAllAccounts(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *AccountRequest) (response *AccountResponse, err error) {
	var (
		body     *AccountResponseBody
		selector *dbstmts.Select[*empty, *accountmodels.AccountRow]
		lg       *slog.Logger = log.With("func", "teamall.getAllTeams", "operation", operation.OperationID)
	)
	lg.Info("starting handler ...")
	timers.Start(operation.OperationID)
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

	// prep result
	timers.Stop(operation.OperationID)
	body = &AccountResponseBody{
		Request:     input,
		Count:       len(selector.Returned),
		Data:        selector.Returned,
		Performance: timers.AllTimers(),
	}
	response = &AccountResponse{Body: body}
	lg.Info("complete.")
	return
}
