package dbselects

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbmodels"
	"opg-reports/report/internal/db/dbstmts"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrTransactionBeginFailed = errors.New("transaction begin failed with error.")
	ErrPreparedStmtFailed     = errors.New("prepared stmt failed with error.")
	ErrMissingResults         = errors.New("error with returned results.")
	ErrFailedTx               = errors.New("error comitting txn.")
	ErrMissingTable           = errors.New("missing table.")
)

// Select creates a transaction to run SQL command within the db. Data is attached to the `.Returned` property
// on `stmt`
func Select[T dbmodels.Model, R dbmodels.Result](ctx context.Context, log *slog.Logger, db *sqlx.DB, stmt *dbstmts.Select[T, R]) (err error) {

	var (
		transaction *sqlx.Tx
		statement   *sqlx.NamedStmt
		lg          *slog.Logger   = log.With("func", "dbselects.Select")
		data        T              = stmt.Data
		returned    []R            = []R{}
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: true, Isolation: sql.LevelDefault}
	)

	lg.Debug("starting ...")
	// start transaction
	transaction, err = db.BeginTxx(ctx, options)
	if err != nil {
		lg.Error("error with transaction begin", "err", err.Error())
		err = errors.Join(ErrTransactionBeginFailed, err)
		return
	}
	// create prepared statement so placeholders are used
	statement, err = transaction.PrepareNamedContext(ctx, stmt.Statement)
	if err != nil {
		lg.Debug("prepared stmt failed", "err", err.Error(), "stmt", stmt.Statement)
		err = errors.Join(ErrPreparedStmtFailed, err)
		if strings.Contains(err.Error(), "no such table") {
			err = errors.Join(ErrMissingTable, err)
		}
		return
	}
	// run the select and attach the result
	err = statement.SelectContext(ctx, &returned, data)
	if err != nil && err != sql.ErrNoRows {
		lg.Error("stmt context failed", "error", err.Error())
		err = errors.Join(ErrMissingResults, err)
		return
	}
	stmt.Returned = returned
	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		lg.Error("failed to commit transaction")
		err = errors.Join(ErrFailedTx, err)
	}

	lg.With("count", len(returned)).Debug("complete")
	return
}
