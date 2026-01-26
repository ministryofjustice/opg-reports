package dbexec

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstatements"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Exec runs a complete statement against the database and returns any error
func Exec(ctx context.Context, log *slog.Logger, db *sqlx.DB, statement dbstatements.Statement) (result sql.Result, err error) {
	var (
		transaction *sqlx.Tx
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}
	)

	log = log.With("package", "db.dbexec", "func", "Exec", "statement", string(statement))
	log.Debug("starting ...")
	// start transaction
	transaction, err = db.BeginTxx(ctx, options)
	if err != nil {
		log.Error("error with transaction begin", "err", err.Error())
		err = errors.Join(ErrTransactionBeginFailed, err)
		return
	}
	// try to execute
	result, err = transaction.ExecContext(ctx, string(statement))
	if err != nil {
		log.Error("error with transaction exec", "err", err.Error())
		err = errors.Join(ErrTransactionExecFailed, err)
		return
	}

	// no error, so commit the transaction
	err = transaction.Commit()
	if err != nil {
		log.Error("error with transaction commit", "err", err.Error())
		err = errors.Join(ErrTransactionCommitFailed, err)
		// rollback
		transaction.Rollback()
		return
	}

	log.Debug("completed.")
	return
}
