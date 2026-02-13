package dbsetup

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbexec"
	"opg-reports/report/internal/db/dbinserts"
	"opg-reports/report/internal/db/dbstmts"

	"github.com/jmoiron/sqlx"
)

var (
	ErrPreFailed    = errors.New("import pre-statement failed with error.")
	ErrImportFailed = errors.New("import failed with error.")
)

// ImportStatements are run combined with the data, the `Pre` attribute
// allows standalone sql (like truncates) to be run before insert
type ImportStatements struct {
	Pre    string
	Insert string
}

func Import[R any, T any](ctx context.Context, log *slog.Logger, db *sqlx.DB, data []T, statements *ImportStatements) (inserted []*dbstmts.Insert[T, R], err error) {
	var lg *slog.Logger = log.With("func", "dbimports.Import", "types", fmt.Sprintf("%T", inserted))

	lg.Debug("starting ...")
	inserted = []*dbstmts.Insert[T, R]{}

	// run any pre statement - generally a truncate
	if statements.Pre != "" {
		lg.Debug("running pre statement ...")
		_, err = dbexec.Exec(ctx, log, db, dbstmts.Statement(statements.Pre))
	}
	// check for errors
	if err != nil {
		lg.Error("pre statement failed.")
		err = errors.Join(ErrPreFailed, err)
		return
	}
	// generate insert statements for the data rows passed
	lg.Debug("generating db insert statements ...")
	for _, row := range data {
		inserted = append(inserted, &dbstmts.Insert[T, R]{
			Statement: statements.Insert,
			Data:      row,
		})
	}
	lg.Debug("running import statements via insert ...")
	err = dbinserts.Insert(ctx, log, db, inserted...)
	if err != nil {
		lg.Error("error with insert.", "err", err.Error())
		err = errors.Join(ErrImportFailed, err)
		return
	}
	lg.Debug("complete.")
	return
}
