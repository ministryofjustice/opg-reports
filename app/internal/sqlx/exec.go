package sqlx

import (
	"context"
	"database/sql"
	"log/slog"
	"opg-reports/app/internal/logx"
)

// Exec provides a more raw approach to running SQL against the database. No binding with structs
// is attempted.
//
// Generally used to run TABLE creation or DELETE commands within migration etc.
//
// Notes: Does not close the connection
func Exec(ctx context.Context, conn Connector, sqlStmt string, args ...any) (result sql.Result, err error) {
	var (
		db *sql.DB
		lg *slog.Logger = logx.Default()
	)
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

	result, err = db.ExecContext(ctx, sqlStmt, args...)
	if err != nil {
		lg.Error("error running sql.",
			"sql", sqlStmt,
			"err", err.Error())
	}

	return
}
