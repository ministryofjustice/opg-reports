package sqlx

import (
	"context"
	"database/sql"
	"log/slog"
	"opg-reports/app/internal/logx"
)

// preparedInsert is used by insert to group the result of Bind
// together so all records can be checked against the sql statement
// before executing.
type preparedInsert struct {
	Sql  string
	Args []interface{}
}

// Insert intends to provide a way to run a sql statement where fields placeholders (`:fieldName`) are merged with values
// from a struct. Intended to be used for INSERT and UPDATE.
//
// All preparted statements are created first, via the `Bind` function and grouped with their respective arguments. This
// means any errors in the bindings will happen before any executions.
//
// SQL is then executed via `ExecContext`. Any error from the SQL will stop all others and return.
//
// Notes:
// Currently runs outside of a transaction, so there is no roll-back.
// Expects db connection to be writable.
// Does not close the db connection
func Insert[T any](ctx context.Context, conn Connector, sqlStmt string, records []T) (results []sql.Result, err error) {
	var (
		db      *sql.DB
		lg      *slog.Logger      = logx.Default()
		prepped []*preparedInsert = []*preparedInsert{}
	)
	results = []sql.Result{}
	// check connection
	if conn.Mode() == READONLY {
		lg.Error("read-only mode connection.")
		err = ErrReadOnlyMode
		return
	}
	// open the db connection
	db, err = conn.Open()
	if err != nil {
		lg.Error("failed to open database connection.", "err", err.Error())
		return
	}
	// loop over all the records to create the set of prepared statements
	for _, entry := range records {
		sqlStr, args, e := Bind(ctx, sqlStmt, entry)
		// if bind fails, return an error
		if e != nil {
			err = e
			lg.Error("failed to bind entry to sql.", "sqlStmt", sqlStmt, "err", e.Error())
			return
		}

		// add a working sql & arg combination to the list
		prepped = append(prepped, &preparedInsert{
			Sql:  sqlStr,
			Args: args,
		})
	}
	// loop over the prepared statements and execute them
	for _, cmd := range prepped {
		res, e := db.ExecContext(ctx, cmd.Sql, cmd.Args...)
		if e != nil {
			err = e
			lg.Error("failed to exectute sql.", "sql", cmd.Sql, "err", e.Error())
			return
		}
		results = append(results, res)
	}

	return
}
