package dbx

import (
	"context"
	"database/sql"
	"log/slog"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/fmtx"
	"opg-reports/report/packages/logger"
	"opg-reports/report/packages/types/interfaces"
)

// Select
func Select[T interfaces.Selectable](ctx context.Context, statement string, filter interfaces.Filterable, opts *args.DB) (results []T, err error) {
	var (
		log       *slog.Logger
		db        *sql.DB
		rows      *sql.Rows
		values    []interface{}
		resultSet Results[T] = Results[T]{}
	)
	ctx, log = logger.Get(ctx)

	db, err = sql.Open(opts.Connection())
	if err != nil {
		log.Error("select: error connecting to database", "err", err.Error())
		return
	}
	defer db.Close()
	// use SprintfNamed to update the sql (resolving slices etc) to have placeholders
	// and return the values of them
	statement, values = fmtx.SprintfNamed(statement, filter.Map(), true)
	// setup row scanning
	rows, err = db.QueryContext(ctx, statement, values...)
	if err != nil {
		log.Debug("select: statement with error", "statement", statement)
		log.Error("select: error in query", "err", err.Error())
		return
	}

	defer rows.Close()
	// now interate over the rows
	for rows.Next() {
		err = resultSet.RowScan(rows)
		if err != nil {
			log.Error("unexpected error with row scanning.", "err", err.Error())
			return
		}
	}
	// return the results only
	results = resultSet.Data()

	return

}
