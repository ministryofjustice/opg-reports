package codebaseimports

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbimports"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codebases/codebasemodels"

	"github.com/jmoiron/sqlx"
)

var ErrTruncateFailed = errors.New("codebase truncate failed with error.")

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO codebases (
	name,
	full_name,
	url
) VALUES (
	:name,
	:full_name,
	:url
)
ON CONFLICT (full_name)
 	DO UPDATE SET
		name=excluded.name,
		url=excluded.url
RETURNING id
;
`

// truncate to remove all entries as code bases may be removed / changed over time
const truncateStmt string = `DELETE FROM codebases;`

// Import uses combines the cost data passed along with the with insert statement defined in this package to
// insert records in to the active database connection.
//
// Table is truncated first, as codebases may have changed over time (deleted / renamed etc).
func Import(ctx context.Context, log *slog.Logger, db *sqlx.DB, data []*codebasemodels.Codebase) (statements []*dbstmts.Insert[*codebasemodels.Codebase, int], err error) {
	var lg *slog.Logger = log.With("func", "codebaseimports.Import")

	lg.Debug("starting ...")
	lg.Debug("trucating table ...")
	// truncate table
	_, err = dbexec.Exec(ctx, log, db, dbstmts.Statement(truncateStmt))
	if err != nil {
		lg.Error("error with truncate.", "err", err.Error())
		err = errors.Join(ErrTruncateFailed, err)
		return
	}

	statements, err = dbimports.Import[int](ctx, log, db, insertStmt, data)
	if err != nil {
		return
	}
	lg.Debug("complete.")
	return
}
