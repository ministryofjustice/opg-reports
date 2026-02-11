package infracostselects

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"

	"github.com/jmoiron/sqlx"
)

const selectAllStmt string = `
SELECT
	id,
	region,
	service,
	date,
	cost,
	account_id
FROM infracosts
ORDER BY
	date,
	CAST(cost as REAL) DESC,
	account_id,
	service DESC
`

// empty is used when no filtering / subsitution on the sql statement, such
// as a `select *` or `select count(*)`
type empty struct{}

func All(ctx context.Context, log *slog.Logger, db *sqlx.DB) (code []*infracostmodels.Cost, err error) {
	var (
		lg       *slog.Logger = log.With("func", "infracostselects.All")
		selector *dbstmts.Select[*empty, *infracostmodels.Cost]
	)
	code = []*infracostmodels.Cost{}
	// setup the select
	selector = &dbstmts.Select[*empty, *infracostmodels.Cost]{
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
