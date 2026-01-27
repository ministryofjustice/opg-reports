package dbinserts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"opg-reports/report/internal/db/dbmodels"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Insert creates a transaction for each statement and will fail if any insert fails.
//
// Return data is populated within the statements directly.
//
// Note: On fail a rollback is triggered as well as error being returned.
func Insert[T dbmodels.Model, R dbmodels.Result](ctx context.Context, log *slog.Logger, db *sqlx.DB, statements ...*dbstatements.InsertStatement[T, R]) (err error) {
	var (
		transaction *sqlx.Tx
		options     *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}
	)

	log = log.With("package", "dbinserts", "func", "Insert")
	// start transaction
	transaction, err = db.BeginTxx(ctx, options)
	if err != nil {
		log.Error("error with transaction begin", "err", err.Error())
		err = errors.Join(ErrTransactionBeginFailed, err)
		return
	}

	log.Debug("generating set of prepared context statements ...")
	for _, stmt := range statements {
		var statement *sqlx.NamedStmt
		var data = stmt.Data

		// check error on prep
		statement, err = transaction.PrepareNamedContext(ctx, stmt.Statement)
		if err != nil {
			log.Error("prepared insert stmt failed", "err", err.Error(), "stmt", stmt.Statement)
			err = errors.Join(ErrPreparedInsertFailed, err)
			return
		}
		// check error on context
		err = statement.GetContext(ctx, &stmt.Returned, data)
		if err != nil && err != sql.ErrNoRows {
			log.Error("stmt context failed", "err", err.Error(), "sql", statement.QueryString)
			err = errors.Join(ErrGetContextFailed, err)
			return
		}
	}

	log.Debug("executing prepared contexts in transaction ...")
	err = transaction.Commit()
	if err != nil {
		log.Error("error with transaction commit", "err", err.Error())
		err = errors.Join(ErrTransactionCommitFailed, err)
		// rollback
		transaction.Rollback()
		return
	}
	// check results are aligned and have values
	err = checkSuccess(statements...)
	if err != nil {
		log.Error("error with missing results", "err", err.Error())
		err = errors.Join(ErrMissingResults, err)
		return
	}

	return
}

// checkSuccess looks at each statments to see if there are any without a
// Returned value, if so it flags an error
func checkSuccess[T dbmodels.Model, R dbmodels.Result](statements ...*dbstatements.InsertStatement[T, R]) (err error) {
	// check if everything was inserted
	for _, stmt := range statements {
		var r = utils.Ptr(stmt.Returned)
		if r == nil {
			return fmt.Errorf("some statments failed to complete")
		}
	}
	return nil
}
