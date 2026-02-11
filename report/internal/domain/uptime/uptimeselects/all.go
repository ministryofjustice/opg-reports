package uptimeselects

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/uptime/uptimemodels"

	"github.com/jmoiron/sqlx"
)

const selectAllStmt string = `
SELECT
	id,
	date,
	average,
	granularity,
	account_id
FROM uptime
ORDER BY date, account_id DESC
`

// empty is used when no filtering / subsitution on the sql statement, such
// as a `select *` or `select count(*)`
type empty struct{}

func All(ctx context.Context, log *slog.Logger, db *sqlx.DB) (code []*uptimemodels.Uptime, err error) {
	var (
		lg       *slog.Logger = log.With("func", "teamselects.All")
		selector *dbstmts.Select[*empty, *uptimemodels.Uptime]
	)
	code = []*uptimemodels.Uptime{}
	// setup the select
	selector = &dbstmts.Select[*empty, *uptimemodels.Uptime]{
		Statement: selectAllStmt,
		Data:      &empty{},
	}

	lg.Debug("starting ...")
	err = dbselects.Select(ctx, log, db, selector)
	if err != nil {
		lg.Error("error with select.", "err", err.Error())
		return
	}
	// setup output
	for _, row := range selector.Returned {
		code = append(code, row)
	}
	lg.Debug("complete.")
	return
}
