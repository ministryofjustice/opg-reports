package codebaseall

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/utils/timers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
)

// fixed values for this endpoint, used by the operation setup for huma
const (
	ENDPOINT      string = `/v1/codebases/all`
	opID          string = `codebases-get-all`
	opSummary     string = `Return all codebases and their codeowners.`
	opDescription string = `Returns a list of all codebases and codeowner data from the database without filtering.`
)

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

// errors
var (
	ErrSelectFailed  = errors.New("select query failed with error.")
	ErrConvertFailed = errors.New("dataset type conversion failed.")
)

// selectStmt is the main sql statment to fetch data from the db
var selectStmt string = `
SELECT
	codebases.id,
	codebases.name,
	codebases.full_name,
	codebases.url,
	json_group_array(
		DISTINCT json_object(
			'name', codeowners.name,
			'team_name', codeowners.team_name
		)
	) filter ( where codeowners.name is not null) as codeowner_list
FROM codebases
LEFT JOIN codeowners on codeowners.codebase_full_name = codebases.full_name
GROUP BY codebases.full_name
ORDER BY
	codeowners.team_name,
	full_name ASC
;`

// Request contains the incoming url and query string data for this endpoint
type CodebaseRequest struct{}

// Response is the handlers data struct passed to a huma api which will then be rendered
type CodebaseResponse struct {
	Body *CodebaseResponseBody
}

// ResponseBody is the response body, containing all data to be returned
type CodebaseResponseBody struct {
	Request     *CodebaseRequest              `json:"request"`
	Data        []*codebasemodels.CodebaseAll `json:"data"`
	Count       int                           `json:"count,omitempty"`
	Performance []*timers.Timer               `json:"performance"`
}

// empty is used as the input data for the select statement
// as there are now filters etc
type empty struct{}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {

	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *CodebaseRequest) (*CodebaseResponse, error) {
		return getAll(ctx, log, db, &operation, input)
	})

}

// getAll
func getAll(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *CodebaseRequest) (response *CodebaseResponse, err error) {
	var (
		body     *CodebaseResponseBody
		selector *dbstmts.Select[*empty, *codebasemodels.CodebaseAll]
		lg       *slog.Logger = log.With("func", "codebaseapis.getAll", "operation", operation.OperationID)
	)
	lg.Info("starting handler ...")
	timers.Start(operation.OperationID)
	// create the statement
	lg.Debug("creating select statement ...")
	selector = &dbstmts.Select[*empty, *codebasemodels.CodebaseAll]{
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
	body = &CodebaseResponseBody{
		Request:     input,
		Count:       len(selector.Returned),
		Data:        selector.Returned,
		Performance: timers.AllTimers(),
	}
	response = &CodebaseResponse{Body: body}
	lg.Info("complete.")
	return
}
