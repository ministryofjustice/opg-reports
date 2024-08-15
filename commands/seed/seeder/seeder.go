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

func Seed(ctx context.Context, dbF string, schemaF string, dataF string, table string, N int) (db *sql.DB, err error) {
	logger.LogSetup()
	var insertFunc insertF = nil
	var trackFunc trackerF = nil

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

	// look for import and tracking functions
	if ik, ok := INSERT_FUNCTIONS[table]; ok {
		insertFunc = ik
	}
	if tk, ok := TRACKER_FUNCTIONS[table]; ok {
		trackFunc = tk
	}

	if len(files) > 0 && insertFunc != nil {
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

	} else if fk, ok := GENERATOR_FUNCTIONS[table]; ok && len(files) == 0 {
		slog.Info("no files, but generator function - creating dummy data")
		if err = fk(ctx, N, db); err != nil {
			return
		}
		// -- track the generation
		trackFunc(ctx, time.Now().UTC(), db)
	}

	return
}

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
