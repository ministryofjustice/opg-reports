package dbselects

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	ErrNamedFailed    = errors.New("error when calling named.")
	ErrInRebindFailed = errors.New("error when calling sqlx in.")
	ErrSelectFailed   = errors.New("error running select.")
)

func Select2[T dbmodels.Model, R dbmodels.Result](ctx context.Context, log *slog.Logger, db *sqlx.DB, stmt *dbstmts.Select[T, R]) (err error) {

	var (
		query  string        // the generated sqlx query statement
		args   []interface{} // the generated arguments used for a query
		lg     *slog.Logger  = log.With("func", "dbselects.Select2", "stmt", fmt.Sprintln(stmt.Statement))
		result []R           = []R{}
	)

	lg.Debug("starting ...")
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

	// run the select
	err = db.Select(&result, query, args...)
	if err != nil {
		lg.Error("error running select", "err", err.Error())
		err = errors.Join(ErrSelectFailed, err)
		return
	}

	stmt.Returned = result
	lg.With("count", len(result)).Debug("complete")

	// query, args, err = sqlx.Named(stmt.Statement, stmt.Data)
	// if err != nil {
	// 	return
	// }
	// query, args, err = sqlx.In(query, args...)
	// if err != nil {
	// 	return
	// }

	// query = db.Rebind(query)
	// fmt.Println("query >>")
	// debugger.Dump(query)

	// // rows, err := db.Query(query, args...)
	// var res []R = []R{}
	// db.Select(&res, query, args...)

	// fmt.Println("query >>")
	// debugger.Dump(query)
	// fmt.Println("args >>")
	// debugger.Dump(args)
	// fmt.Println("res >>")
	// debugger.Dump(res)
	// fmt.Println("err >>")
	// debugger.Dump(err)
	// fmt.Println("====")
	return

}

// Select creates a transaction to run SQL command within the db. Data is attached to the `.Returned` property
// on `stmt`
func Select[T dbmodels.Model, R dbmodels.Result](ctx context.Context, log *slog.Logger, db *sqlx.DB, stmt *dbstmts.Select[T, R]) (err error) {

	var (
		transaction *sqlx.Tx
		statement   *sqlx.NamedStmt
		lg          *slog.Logger   = log.With("func", "dbselects.Select")
		args        T              = stmt.Data
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
		lg.Warn("prepared stmt failed", "err", err.Error(), "stmt", stmt.Statement, "data", stmt.Data)
		err = errors.Join(ErrPreparedStmtFailed, err)
		if strings.Contains(err.Error(), "no such table") {
			err = errors.Join(ErrMissingTable, err)
		}
		return
	}
	// run the select and attach the result
	err = statement.SelectContext(ctx, &returned, args)
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
