// Package seeder main function is to generate sqlite3 databases for use with the api.
//
// It uses schemas from sqlc to create the tagles within the database and will either seed
// those tables with data from a series of existing files or by generating new|fake data.
//
// Seed is the main function to call and will create the database and insert
// data.
//
// Used by:
//   - ./commands/seed (which in turn is called within Dockerfile)
//   - various test files to create dummy data to test against the api
package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/exists"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

type generatorF func(ctx context.Context, num int, db *sql.DB) error
type insertF func(ctx context.Context, fileContent []byte, db *sql.DB) error
type trackerF func(ctx context.Context, ts time.Time, db *sql.DB) error

// Seed tries to create a database at the filepath set (`dbF`) and return the db pointer.
// This pointer is *NOT* closed, it will need to be handled outside of this function
//
// If there is already a file at `dbF` location, this will exit without error, but will *NOT* create a new version.
//
// A database connection is created (and returned) using the consts.SQL_CONNECTION_PARAMS. The schema file (`schemaF`)
// is then read and passed into the database to execute and generate the empty tables and indexes. This schema is part
// of this projects use of sqlc
//
// An insert and tracker function are then checked for, comparing the table name (`table`) against the known list. These are
// required to be used, so the function will error if a match is not found. See generators.go, insertors.go and trackers.go
// for list of current versions.
//
// The file pattern (`dataF`) is then used with Glob to find matching files. If they are found, they are iterated over
// and their contents passed into the insertor function - this function will handle marshaling / coversion and the db
// record creation
//
// If no files are found but there is a generator function, this will be called instead to place dummy data into the
// database. The amount of records created is controlled by the `N` parameter
func Seed(ctx context.Context, dbF string, schemaF string, dataF string, table string, N int) (db *sql.DB, err error) {
	logger.LogSetup()
	var (
		ok         bool       = false
		insertFunc insertF    = nil
		trackFunc  trackerF   = nil
		genFunc    generatorF = nil
	)

	// -- debug info
	slog.Debug("starting to seed",
		slog.String("db", dbF),
		slog.String("schema", schemaF),
		slog.String("dataFiles", dataF),
		slog.String("table", table),
		slog.Int("n", N))

	if dbF == "" || schemaF == "" || dataF == "" || table == "" || N == 0 {
		err = fmt.Errorf("missing required arguments")
		return
	}
	// if the database exists, ignore
	dbExists := exists.FileOrDir(dbF)
	if dbExists {
		err = fmt.Errorf("database already exists: %s", dbF)
		return
	}
	// -- generate connection
	db, err = DB(dbF)
	if err != nil {
		return
	}
	// if the schema does not exist, error
	schemaExists := exists.FileOrDir(schemaF)
	if !schemaExists {
		err = fmt.Errorf("schema does not exist: %s", schemaF)
		return
	}
	// -- load schema into db
	err = SchemaLoad(db, schemaF)
	if err != nil {
		return
	}

	// -- now look for files matching the source pattern
	files := []string{}
	files, err = filepath.Glob(dataF)
	if err != nil {
		slog.Error("error with pattern matching", slog.String("data", dataF), slog.String("err", err.Error()))
	}

	slog.Debug("data files found", slog.Int("count", len(files)))

	// look for import, generator and tracking functions
	insertFunc, ok = haveFuncforTable(table, INSERT_FUNCTIONS)
	if !ok {
		err = fmt.Errorf("no seeder insertion function found for table: %s", table)
		return
	}
	genFunc, ok = haveFuncforTable(table, GENERATOR_FUNCTIONS)
	if !ok {
		err = fmt.Errorf("no generator function found for table: %s", table)
		return
	}
	trackFunc, ok = haveFuncforTable(table, TRACKER_FUNCTIONS)
	if !ok {
		err = fmt.Errorf("no tracker function found for table: %s", table)
		return
	}

	// if we have files and a function, try and insert the content of them into the database
	// otherwise, insert dummy data via the generator function
	// - both will also call the track
	if len(files) > 0 {
		times := []time.Time{}
		for _, file := range files {
			var ts *time.Time = nil
			slog.Debug("found data file", slog.String("file", file))
			err, ts = InsertFile(ctx, file, insertFunc, db)
			if err != nil {
				return
			} else if ts != nil {
				times = append(times, *ts)
			}
		}
		if len(times) > 0 {
			max := dates.MaxTime(times)
			trackFunc(ctx, max, db)
		}

	} else if genFunc != nil {
		slog.Info("no files, but generator function - creating dummy data")
		err = genFunc(ctx, N, db)
		if err != nil {
			return
		}
		// -- track the generation
		trackFunc(ctx, time.Now().UTC(), db)
	}

	return
}

// haveFuncforTable is a mall helper to check we have a function for this table
func haveFuncforTable[T insertF | generatorF | trackerF](table string, set map[string]T) (f T, found bool) {
	f, found = set[table]
	return
}

// InsertFile checks the file exists and reads its content into the insertFunc along with the db pointer.
//
// Small wrapper that also find the modification time of the file. This time is then used to track the age of the
// data being used
func InsertFile(ctx context.Context, file string, insertFunc insertF, db *sql.DB) (err error, ts *time.Time) {
	var ct time.Time
	ts = nil
	if !exists.FileOrDir(file) {
		return
	}
	ct, err = dates.FileCreationTime(file)
	if err != nil {
		return
	}
	ts = &ct

	slog.Debug("insertng content from file",
		slog.String("file", file), slog.String("created", ct.String()))
	content, err := os.ReadFile(file)
	if err != nil {
		return
	}
	// try inserting
	err = insertFunc(ctx, content, db)
	if err != nil {
		return
	}
	defer os.Remove(file)
	return
}

// SchemaLoad generates a schema from a sqlc file whose filepath is passed
func SchemaLoad(db *sql.DB, schemaFile string) (err error) {
	schemaDir := filepath.Dir(schemaFile)
	os.MkdirAll(schemaDir, os.ModePerm)
	// -- read file
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		return
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		return
	}
	return
}

// create a db
func DB(dbFile string) (db *sql.DB, err error) {
	dbDir := filepath.Dir(dbFile)
	os.MkdirAll(dbDir, os.ModePerm)
	// -- create a stub database file
	os.WriteFile(dbFile, []byte(""), os.ModePerm)

	// -- connect to the database and then import schema
	conn := consts.SQL_CONNECTION_PARAMS
	db, err = sql.Open("sqlite3", dbFile+conn)
	if err != nil && db != nil {
		db.Close()
	}
	return
}
