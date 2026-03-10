package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/fmtx"
	"opg-reports/report/packages/logger"
	"opg-reports/report/packages/types/interfaces"

	_ "github.com/mattn/go-sqlite3"
)

// Insert writes records to the database using the sql statement as a template for each insert.
//
// Converts from the []T into generic map[string]interface{} and then uses the `SprintfNamed`
// function to create the fully expanded sql and the correctly ordered arguments.
func Insert[T interfaces.Insertable](ctx context.Context, statement string, records []T, opts *args.DB) (err error) {
	var (
		log *slog.Logger
		db  *sql.DB
	)
	ctx, log = logger.Get(ctx)
	log.Debug("inserting records ...", "type", fmt.Sprintf("%T", records))

	db, err = sql.Open(opts.Connection())
	if err != nil {
		log.Error("insert: error connecting to database", "err", err.Error())
		return
	}
	defer db.Close()

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
	log.Debug("insert sql completed.")
	return
}
