package codeownerimports

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"

	"github.com/jmoiron/sqlx"
)

var (
	ErrImportFailed   = errors.New("codeowners import failed with error.")
	ErrTruncateFailed = errors.New("codeowners truncate failed with error.")
)

// insertStmt used to insert records
const insertStmt string = `
INSERT INTO codeowners (
	name,
	codebase_full_name,
	team_name
) VALUES (
 	:name,
	:codebase_full_name,
	:team_name
)
ON CONFLICT (name,codebase_full_name,team_name)
 	DO UPDATE SET
		team_name=excluded.team_name,
		name=excluded.name,
		team_name=excluded.team_name
RETURNING id
;
`

// truncate to remove all entries as code bases may be removed / changed over time
const truncateStmt string = `DELETE FROM codeowners;`

// Import uses combines the cost data passed along with the with insert statement defined in this package to
// insert records in to the active database connection.
//
// Table is truncated first, as codebases may have changed over time (deleted / renamed etc).
func Import(ctx context.Context, log *slog.Logger, db *sqlx.DB, data []*codeownermodels.Codeowner) (statements []*dbstmts.Insert[*codeownermodels.Codeowner, int], err error) {
	var lg *slog.Logger = log.With("func", "domain.codeowners.codeownerimports.Import")

	statements = []*dbstmts.Insert[*codeownermodels.Codeowner, int]{}
	lg.Debug("starting ...")
	lg.Debug("generating db insert statements ...")
	// generate all of the insert statements from the data passed
	for _, row := range data {
		statements = append(statements, &dbstmts.Insert[*codeownermodels.Codeowner, int]{
			Statement: insertStmt,
			Data:      row,
		})
	}
	// truncate the table
	_, err = dbexec.Exec(ctx, log, db, dbstmts.Statement(truncateStmt))
	if err != nil {
		lg.Error("error with truncate", "err", err.Error())
		err = errors.Join(ErrTruncateFailed, err)
		return
	}

	// run inserts
	lg.Debug("running import statements via insert")
	err = dbinserts.Insert(ctx, log, db, statements...)
	if err != nil {
		lg.Error("error with insert", "err", err.Error())
		err = errors.Join(ErrImportFailed, err)
		return
	}
	lg.Debug("complete.")
	return
}
