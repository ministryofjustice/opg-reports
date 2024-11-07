package datastore

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
	"github.com/ministryofjustice/opg-reports/pkg/record"
)

// Create will uses ExecContext to run the slice of createStatements passed.
// Used to create table, add indexes and so on in sequence
// Any errors in this process will trigger a panic and exit
func Create(ctx context.Context, db *sqlx.DB, create []CreateStatement) {
	slog.Debug("[datastore.Create] ")
	for _, stmt := range create {
		if _, err := db.ExecContext(ctx, string(stmt)); err != nil {
			slog.Error("error in create", slog.String("err", err.Error()))
			panic(err)
		}

	}
}

// CreateNewDB will create a new DB file and then
// try to run table and index creates against it.
// Returns the db, a bool to say if it was new and any errors
func CreateNewDB(ctx context.Context, dbFilepath string, creates []CreateStatement) (db *sqlx.DB, isNew bool, err error) {

	db, isNew, err = NewDB(ctx, Sqlite, dbFilepath)
	if err == nil && len(creates) > 0 {
		Create(ctx, db, creates)
	}

	return
}

// Setup will ensure a database with records exists in the filepath requested.
// If there is no database at that location a new sqlite database will
// be created and populated with series of dummy data - helpful for local testing.
func Setup[T record.Record](ctx context.Context, dbFilepath string, insertStmt InsertStatement, creates []CreateStatement, seed bool, n int) {

	var err error
	var db *sqlx.DB
	var isNew bool = false
	// add custom fakers
	exfaker.AddProviders()

	db, isNew, err = CreateNewDB(ctx, dbFilepath, creates)
	defer db.Close()

	if err != nil {
		panic(err)
	}

	if seed && isNew {
		faked := exfaker.Many[T](n)
		_, err = InsertMany(ctx, db, insertStmt, faked)
	}
	if err != nil {
		panic(err)
	}

}
