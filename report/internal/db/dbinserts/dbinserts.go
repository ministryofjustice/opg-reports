package dbinserts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbmodels"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/utils/ptr"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrTransactionBeginFailed  = errors.New("transaction begin failed with error.")
	ErrPreparedInsertFailed    = errors.New("prepared insert stmt failed with error.")
	ErrGetContextFailed        = errors.New("stmt context failed with error.")
	ErrTransactionExecFailed   = errors.New("transaction insert failed with error.")
	ErrTransactionCommitFailed = errors.New("transaction commit failed with error.")
	ErrMissingResults          = errors.New("error with returned results.")
)

// Insert creates a transaction for each statement and will fail if any insert fails.
//
// Return data is populated within the statements directly.
//
// Note: On fail a rollback is triggered as well as error being returned.
func Insert[T dbmodels.Model, R dbmodels.Result](ctx context.Context, log *slog.Logger, db *sqlx.DB, statements ...*dbstmts.Insert[T, R]) (err error) {
	var (
		transaction *sqlx.Tx
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}
		lg          *slog.Logger   = log.With("func", "dbinserts.Insert")
	)
	lg.Debug("starting ...")
	// start transaction
	transaction, err = db.BeginTxx(ctx, options)
	if err != nil {
		lg.Error("error with transaction begin", "err", err.Error())
		err = errors.Join(ErrTransactionBeginFailed, err)
		return
	}

	lg.Debug("generating set of prepared context statements ...")
	for _, stmt := range statements {
		var statement *sqlx.NamedStmt
		var data = stmt.Data

		// check error on prep
		statement, err = transaction.PrepareNamedContext(ctx, stmt.Statement)
		if err != nil {
			lg.Error("prepared insert stmt failed", "err", err.Error(), "stmt", stmt.Statement)
			err = errors.Join(ErrPreparedInsertFailed, err)
			return
		}
		// check error on context
		err = statement.GetContext(ctx, &stmt.Returned, data)
		if err != nil && err != sql.ErrNoRows {
			lg.Error("stmt context failed", "err", err.Error(), "sql", statement.QueryString)
			err = errors.Join(ErrGetContextFailed, err)
			return
		}
	}

	lg.Debug("executing prepared contexts in transaction ...")
	err = transaction.Commit()
	if err != nil {
		lg.Error("error with transaction commit", "err", err.Error())
		err = errors.Join(ErrTransactionCommitFailed, err)
		// rollback
		transaction.Rollback()
		return
	}
	// check results are aligned and have values
	err = checkSuccess(statements...)
	if err != nil {
		lg.Error("error with missing results", "err", err.Error())
		err = errors.Join(ErrMissingResults, err)
		return
	}
	lg.Debug("complete.")
	return
}

// checkSuccess looks at each statments to see if there are any without a
// Returned value, if so it flags an error
func checkSuccess[T dbmodels.Model, R dbmodels.Result](statements ...*dbstmts.Insert[T, R]) (err error) {
	// check if everything was inserted
	for _, stmt := range statements {
		var r = ptr.Ptr(stmt.Returned)
		if r == nil {
			return fmt.Errorf("some statments failed to complete")
		}
	}
	return nil
}
