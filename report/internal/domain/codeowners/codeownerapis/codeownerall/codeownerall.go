package codeownerall

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/utils/timers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/codeowner/all`
	opID          string = `codeowner-get-all`
	opSummary     string = `Return all codeowner data`
	opDescription string = `Returns a list of all codeowners and codebase data`
)

// operation describes what this endpoint is doing
var operation = huma.Operation{
	Method:        http.MethodGet,
	DefaultStatus: http.StatusOK,
	Path:          ENDPOINT,
	Summary:       opSummary,
	Description:   opDescription,
	OperationID:   opID,
	Tags:          []string{"codeowner"},
}

// errors
var (
	ErrSelectFailed  = errors.New("select query failed with error.")
	ErrConvertFailed = errors.New("dataset type conversion failed.")
)

// selectStmt is the main sql statment to fetch data from the db
var selectStmt string = `
SELECT
	codeowners.name,
	codeowners.team_name,
	codeowners.codebase_full_name,
	codebases.url
FROM codeowners
LEFT JOIN codebases on codeowners.codebase_full_name = codebases.full_name
ORDER BY
	codeowners.team_name,
	codeowners.codebase_full_name ASC
;`

// Request contains the incoming url and query string data for this endpoint
type CodeownerRequest struct{}

// Response is the handlers data struct passed to a huma api which will then be rendered
type CodeownerResponse struct {
	Body *CodeownerResponseBody
}

// ResponseBody is the response body, containing all data to be returned
type CodeownerResponseBody struct {
	Request     *CodeownerRequest               `json:"request"`
	Data        []*codeownermodels.CodeownerAll `json:"data"`
	Count       int                             `json:"count,omitempty"`
	Performance []*timers.Timer                 `json:"performance"`
}

// empty is used as the input data for the select statement
// as there are now filters etc
type empty struct{}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *CodeownerRequest) (*CodeownerResponse, error) {
		return getAll(ctx, log, db, &operation, input)
	})

}

// getAll
func getAll(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *CodeownerRequest) (response *CodeownerResponse, err error) {
	var (
		body     *CodeownerResponseBody
		selector *dbstmts.Select[*empty, *codeownermodels.CodeownerAll]
		lg       *slog.Logger = log.With("func", "codeowners.getAll", "operation", operation.OperationID)
	)
	lg.Info("starting handler ...")
	timers.Start(operation.OperationID)
	// create the statement
	lg.Debug("creating select statement ...")
	selector = &dbstmts.Select[*empty, *codeownermodels.CodeownerAll]{
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
	body = &CodeownerResponseBody{
		Request:     input,
		Count:       len(selector.Returned),
		Data:        selector.Returned,
		Performance: timers.AllTimers(),
	}
	response = &CodeownerResponse{Body: body}
	lg.Info("complete.")
	return
}
