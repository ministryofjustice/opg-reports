package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/exists"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

func Seed(ctx context.Context, dbF string, schemaF string, dataF string, table string, N int) (db *sql.DB, err error) {
	logger.LogSetup()
	slog.Info("starting to seed ",
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

	slog.Info("how many data files found?", slog.Int("count", len(files)))

	// -- if we didnt find any files, then we try to make some fake ones
	if fk, ok := fakes[table]; ok && len(files) == 0 {
		slog.Info("no files, but generator function - creating dummy data")
		if err = fk(ctx, N, db); err != nil {
			return
		}
	}

	var insertFunc insertF = nil
	// import files, assume list of
	if ik, ok := inserts[table]; ok {
		insertFunc = ik
	}
	for _, file := range files {
		slog.Info("data file", slog.String("file", file))
		if insertFunc != nil && exists.FileOrDir(file) {
			slog.Info("file exists and found insert function", slog.String("file", file))
			content, err := os.ReadFile(file)
			if err != nil {
				db.Close()
				return nil, err
			}
			// try inserting
			err = insertFunc(ctx, content, db)
			if err != nil {
				db.Close()
				return nil, err
			} else {
				defer os.Remove(file)
			}
		}
	}

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
