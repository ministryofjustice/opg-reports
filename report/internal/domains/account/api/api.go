package api

import (
	"context"
	"database/sql"
	"opg-reports/report/packages/dbx"
	"opg-reports/report/packages/httpx"
	"opg-reports/report/packages/slogx"
	"strings"
)

// label is used as the key in the response data and
// a way to track the handler that was called in the
// log output.
const label string = `get-accounts`

// selectStatement
//
// `{WHERE}` is replaced with real value or removed
const selectStatement string = `
SELECT
	id,
	name,
	label,
	environment,
	team_name as team
FROM accounts
WHERE
	{WHERE}
	id IS NOT NULL
ORDER BY
	team_name,
	name,
	environment ASC
;`

func getSelect(filters *httpx.Filter) (stmt string) {
	var filter = ""
	if filters.Team != "" {
		filter = `team_name = :team AND`
	}
	stmt = strings.ReplaceAll(selectStatement, `{WHERE}`, filter)
	return
}

// Accounts fetches all the account records and optionally filters
// based on the team name passed along via the request
func Accounts(ctx context.Context, m httpx.Mux, r httpx.FitleredRequest, cfg httpx.MuxConfig, response *httpx.ResponseContent) {
	var (
		err        error
		db         *sql.DB = cfg.Connection()
		records            = []*Account{}
		filter             = r.Filter()
		selectStmt         = getSelect(filter)
		log                = slogx.FromContext(ctx)
	)

	log.Info(ctx, "fetching data", "label", label)
	records, err = dbx.Select[*Account](ctx, selectStmt, filter, db)
	if err != nil {
		log.Error(ctx, "error getting data", "err", err.Error())
	}
	response.Data[label] = &Result{
		Accounts: records,
	}
	log.Info(ctx, "data fetch complete.", "label", label, "count", len(records))
}
