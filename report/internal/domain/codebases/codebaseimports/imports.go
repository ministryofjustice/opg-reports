package codebaseimports

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/codebases/codebasemodels"

	"github.com/jmoiron/sqlx"
)

var ErrImportFailed = errors.New("codebase import failed with error")

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

// Import uses combines the cost data passed along with the with insert statement defined in this package to
// insert records in to the active database connection.
//
// `data` is presumed to come from the account.GetAwsAccountData
func Import(ctx context.Context, log *slog.Logger, db *sqlx.DB, data []*codebasemodels.Codebase) (statements []*dbstatements.InsertStatement[*codebasemodels.Codebase, int], err error) {

	statements = []*dbstatements.InsertStatement[*codebasemodels.Codebase, int]{}
	log = log.With("package", "codebases", "func", "Import")

	log.Debug("starting ...")
	log.Debug("generating db insert statements ...")
	// generate all of the insert statements from the data passed
	for _, row := range data {
		statements = append(statements, &dbstatements.InsertStatement[*codebasemodels.Codebase, int]{
			Statement: insertStmt,
			Data:      row,
		})
	}
	// run inserts
	log.Debug("running import statements via insert")
	err = dbinserts.Insert(ctx, log, db, statements...)
	if err != nil {
		log.Error("error with insert", "err", err.Error())
		err = errors.Join(ErrImportFailed, err)
		return
	}
	log.Debug("complete.")
	return
}
