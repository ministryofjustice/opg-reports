package codebaseselects

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codebases/codebasemodels"

	"github.com/jmoiron/sqlx"
)

const selectAllStmt string = `
SELECT
	id,
	name,
	full_name,
	url
FROM codebases
ORDER BY name ASC
`

// empty is used when no filtering / subsitution on the sql statement, such
// as a `select *` or `select count(*)`
type empty struct{}

func All(ctx context.Context, log *slog.Logger, db *sqlx.DB) (code []*codebasemodels.Codebase, err error) {
	var (
		lg       *slog.Logger = log.With("func", "domain.codebases.codebaseselects.All")
		selector *dbstmts.Select[*empty, *codebasemodels.Codebase]
	)
	code = []*codebasemodels.Codebase{}
	// setup the select
	selector = &dbstmts.Select[*empty, *codebasemodels.Codebase]{
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
