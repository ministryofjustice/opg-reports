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
	ErrFailedTx               = errors.New("error comitting txn.")
	ErrTransactionBeginFailed = errors.New("transaction begin failed with error.")
	ErrNamedFailed            = errors.New("error when calling named.")
	ErrInRebindFailed         = errors.New("error when calling sqlx in.")
	ErrSelectFailed           = errors.New("error running select.")
)

// Select creates a transaction to run SQL command within the db. Data is attached to the `.Returned` property
// on `stmt`
//
// If the sql statment contains an IN (` in (`) then extra steps are taken to rebind vars to allow this to work
func Select[T dbmodels.Model, R dbmodels.Result](ctx context.Context, log *slog.Logger, db *sqlx.DB, stmt *dbstmts.Select[T, R]) (err error) {

	var (
		transaction *sqlx.Tx
		query       string                 // the generated sqlx query statement
		args        []interface{}          // the generated arguments used for a query
		result      []R            = []R{} // the result data set
		lg          *slog.Logger   = log.With("func", "dbselects.Select")
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: true, Isolation: sql.LevelDefault}
	)

	lg.With("stmt", stmt.Statement).Debug("starting ...")

	// generate the sql statement and args from passed data
	query, args, err = sqlx.Named(stmt.Statement, stmt.Data)
	if err != nil {
		lg.Error("error when calling named", "err", err.Error())
		err = errors.Join(ErrNamedFailed, err)
		return
	}

	// if there is an IN, handle the args and rebind
	if strings.Contains(strings.ToLower(query), " in (") {
		lg.Debug("sql contains an IN element, rebinding arguments")
		// rebind
		query, args, err = sqlx.In(query, args...)
		if err != nil {
			lg.Error("error calling sqlx.In", "err", err.Error())
			err = errors.Join(ErrInRebindFailed, err)
			return
		}
		query = db.Rebind(query)
	}

	// wrap the select in a transaction
	transaction, err = db.BeginTxx(ctx, options)
	if err != nil {
		lg.Error("error with transaction begin", "err", err.Error())
		err = errors.Join(ErrTransactionBeginFailed, err)
		return
	}

	// run the select
	err = transaction.SelectContext(ctx, &result, query, args...)
	if err != nil {
		lg.Error("error running select", "err", err.Error())
		err = errors.Join(ErrSelectFailed, err)
		return
	}
	stmt.Returned = result

	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		lg.Error("failed to commit transaction")
		err = errors.Join(ErrFailedTx, err)
	}

	lg.With("count", len(result)).Debug("complete")
	return

}
