package dbimports

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbmodels"
	"opg-reports/report/internal/db/dbstmts"

	"github.com/jmoiron/sqlx"
)

var ErrImportFailed = errors.New("import failed with error.")

// Import uses combines the data passed along with the with insert statement defined in this package to
// insert records in to the active database connection.
func Import[R dbmodels.Result, T dbmodels.Model](ctx context.Context, log *slog.Logger, db *sqlx.DB, insertStmt string, data []T) (statements []*dbstmts.Insert[T, R], err error) {
	var lg *slog.Logger = log.With("func", "dbimports.Import")

	statements = []*dbstmts.Insert[T, R]{}
	lg.Debug("starting ...")
	lg.Debug("generating db insert statements ...")
	// generate all of the insert statements from the data passed
	for _, row := range data {
		statements = append(statements, &dbstmts.Insert[T, R]{
			Statement: insertStmt,
			Data:      row,
		})
	}
	// run inserts
	lg.Debug("running import statements via insert ...")
	err = dbinserts.Insert(ctx, log, db, statements...)
	if err != nil {
		lg.Error("error with insert.", "err", err.Error())
		err = errors.Join(ErrImportFailed, err)
		return
	}
	log.Debug("complete.")
	return
}
