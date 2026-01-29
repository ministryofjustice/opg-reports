package codeownerselects

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"

	"github.com/jmoiron/sqlx"
)

const selectAllStmt string = `
SELECT
	id,
	name,
	codebase_full_name,
	team_name
FROM codeowners
ORDER BY codebase_full_name, name ASC
`

// empty is used when no filtering / subsitution on the sql statement, such
// as a `select *` or `select count(*)`
type empty struct{}

func All(ctx context.Context, log *slog.Logger, db *sqlx.DB) (code []*codeownermodels.Codeowner, err error) {
	var (
		lg       *slog.Logger = log.With("func", "domain.codeowners.codeownerselects.All")
		selector *dbstatements.SelectStatement[*empty, *codeownermodels.Codeowner]
	)
	code = []*codeownermodels.Codeowner{}
	// setup the select
	selector = &dbstatements.SelectStatement[*empty, *codeownermodels.Codeowner]{
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
