package dbx

import (
	"context"
	"database/sql"
	"opg-reports/report/packages/fmtx"
	"opg-reports/report/packages/slogx"
)

// Select
func Select[T Selectable](ctx context.Context, statement string, filter Filterable, db *sql.DB) (results []T, err error) {
	var (
		rows      *sql.Rows
		values    []interface{}
		resultSet Results[T] = Results[T]{}
		log                  = slogx.FromContext(ctx)
	)
	// use SprintfNamed to update the sql (resolving slices etc) to have placeholders
	// and return the values of them
	statement, values = fmtx.SprintfNamed(statement, filter.Map(), true)
	// setup row scanning
	rows, err = db.QueryContext(ctx, statement, values...)
	if err != nil {
		log.Debug(ctx, "select: statement with error", "statement", statement)
		log.Error(ctx, "select: error in query", "err", err.Error())
		return
	}

	defer rows.Close()
	// now interate over the rows
	for rows.Next() {
		err = resultSet.RowScan(rows)
		if err != nil {
			log.Error(ctx, "unexpected error with row scanning.", "err", err.Error())
			return
		}
	}
	// return the results only
	results = resultSet.Data()

	return

}
