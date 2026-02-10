package codeownerforteam

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
	ENDPOINT      string = `/v1/codeowner/team/{team}`
	opID          string = `codeowner-get-for-team`
	opSummary     string = `Return all codeowner data for the team requested`
	opDescription string = `Returns a list of all codeowners and codebase data for the requestsed team`
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
WHERE
	codeowners.team_name = :team
ORDER BY
	codeowners.team_name,
	codeowners.codebase_full_name ASC
;`

// Request contains the incoming url and query string data for this endpoint
type CodeownerForTeamRequest struct {
	Team string `json:"team,omitempty" path:"team" doc:"Lowercase version of the team name" example:"sirius" required:"true"`
}

// Response is the handlers data struct passed to a huma api which will then be rendered
type CodeownerForTeamResponse struct {
	Body *CodeownerForTeamResponseBody
}

// ResponseBody is the response body, containing all data to be returned
type CodeownerForTeamResponseBody struct {
	Request     *CodeownerForTeamRequest            `json:"request"`
	Data        []*codeownermodels.CodeownerForTeam `json:"data"`
	Count       int                                 `json:"count,omitempty"`
	Performance []*timers.Timer                     `json:"performance"`
}

// filter is used to provide the filters on the sql statement (`:x` notation) and generally
// populated from a matching request input variable.
type filter struct {
	Team string `json:"team" db:"team"`
}

// Register attachs the local handler to the huma api allows way to pass along the configured logger, db etc
func Register(ctx context.Context, log *slog.Logger, db *sqlx.DB, humaapi huma.API) {
	// input is an empty struct as
	huma.Register(humaapi, operation, func(ctx context.Context, input *CodeownerForTeamRequest) (*CodeownerForTeamResponse, error) {
		return getForTeam(ctx, log, db, &operation, input)
	})

}

// getForTeam
func getForTeam(ctx context.Context, log *slog.Logger, db *sqlx.DB, operation *huma.Operation, input *CodeownerForTeamRequest) (response *CodeownerForTeamResponse, err error) {
	var (
		body     *CodeownerForTeamResponseBody
		selector *dbstmts.Select[*filter, *codeownermodels.CodeownerForTeam]
		lg       *slog.Logger = log.With("func", "codeowners.getForTeam", "operation", operation.OperationID, "in", input)
	)
	lg.Info("starting handler ...")
	timers.Start(operation.OperationID)
	// create the statement
	lg.Debug("creating select statement ...")
	selector = &dbstmts.Select[*filter, *codeownermodels.CodeownerForTeam]{
		Statement: selectStmt,
		Data:      &filter{Team: input.Team},
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
	body = &CodeownerForTeamResponseBody{
		Request:     input,
		Count:       len(selector.Returned),
		Data:        selector.Returned,
		Performance: timers.AllTimers(),
	}
	response = &CodeownerForTeamResponse{Body: body}
	lg.Info("complete.")
	return
}
