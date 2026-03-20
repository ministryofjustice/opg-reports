package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"opg-reports/report/packages/fmtx"
	"opg-reports/report/packages/slogx"

	_ "github.com/mattn/go-sqlite3"
)

// Insert writes records to the database using the sql statement as a template for each insert.
//
// Converts from the []T into generic map[string]interface{} and then uses the `SprintfNamed`
// function to create the fully expanded sql and the correctly ordered arguments.
func Insert[T Insertable](ctx context.Context, statement string, records []T, db *sql.DB) (err error) {
	var (
		log = slogx.FromContext(ctx)
	)
	log.Debug(ctx, "inserting records ...", "type", fmt.Sprintf("%T", records))

	// loop over all records and run the insert on each
	// 	- convert from record into map used .Map()
	// 	- uses SprintfNamed to resolve the sql / args
	for _, model := range records {
		var row = model.Map()
		// get the sql formatted statement using SprintfNamed
		bound, args := fmtx.SprintfNamed(statement, row, true)
		// run the sql
		_, err = db.ExecContext(ctx, bound, args...)
		if err != nil {
			return
		}

	}
	log.Debug(ctx, "insert sql complete.")
	return
}
