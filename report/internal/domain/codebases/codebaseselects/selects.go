package codebaseselects

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbselects"
	"opg-reports/report/internal/db/dbstatements"
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
		selector *dbstatements.SelectStatement[*empty, *codebasemodels.Codebase]
	)
	code = []*codebasemodels.Codebase{}
	// setup the select
	selector = &dbstatements.SelectStatement[*empty, *codebasemodels.Codebase]{
		Statement: selectAllStmt,
		Data:      &empty{},
	}
	log = log.With("package", "codebaseselects", "func", "All")
	log.Debug("starting ...")

	err = dbselects.Select(ctx, log, db, selector)
	if err != nil {
		log.Error("error with select", "err", err.Error())
		return
	}

	for _, row := range selector.Returned {
		code = append(code, row)
	}

	return
}
