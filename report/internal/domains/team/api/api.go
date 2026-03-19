package api

import (
	"context"
	"database/sql"
	"opg-reports/report/packages/dbx"
	"opg-reports/report/packages/httpx"
	"opg-reports/report/packages/slogx"
)

// label is used as the key in the response data and
// a way to track the handler that was called in the
// log output.
const label string = `get-teams`

// selectForNavigationStatement returns most teams, skipping
// org & legacy which we dont want in the navigation
const selectForNavigationStatement string = `
SELECT
	name
FROM teams
WHERE
	name != 'org'
	AND name != 'legacy'
ORDER BY
	name ASC
;`

// TeamsForNavigation fetches all the teams that we want to
// be shown within the front navigation, so it skips the org
// and legacy values
func TeamsForNavigation(ctx context.Context, m httpx.Mux, r httpx.FitleredRequest, cfg httpx.MuxConfig, response *httpx.ResponseContent) {
	var (
		err        error
		db         *sql.DB = cfg.Connection()
		records            = []*Team{}
		filter             = &httpx.Filter{} // empty filter, otherwise r.Filter()
		selectStmt         = selectForNavigationStatement
		log                = slogx.FromContext(ctx)
	)

	log.Info(ctx, "fetching data", "label", label)
	records, err = dbx.Select[*Team](ctx, selectStmt, filter, db)
	if err != nil {
		log.Error(ctx, "error getting data", "err", err.Error())
	}
	response.Data[label] = &Result{
		Teams: records,
	}
	log.Info(ctx, "data fetch complete.", "label", label, "count", len(records))
}
