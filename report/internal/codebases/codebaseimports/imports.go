package codebaseimports

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/codebases/codebasemodels"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstatements"

	"github.com/jmoiron/sqlx"
)

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
