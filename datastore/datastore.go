package datastore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type dbVariation struct {
	YearFormat         string
	YearMonthFormat    string
	YearMonthDayFormat string
}

var (
	Sqlite *dbVariation = &dbVariation{YearFormat: "%Y", YearMonthFormat: "%Y-%m", YearMonthDayFormat: "%Y-%m-%d"}
)

const (
	connectionParams string = "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000"
	driverName       string = "sqlite3"
)

// Generic returns the in value as an interface instead of its original type
// to allow single types in responses etc
func Generic(in interface{}, e error) (result interface{}, err error) {
	err = e
	if err == nil {
		result = in
	}
	return
}

// Concreate casts the interface passed to the
// type T if possible
// Generally paired with the output for Generic
func Concreate[T any](in interface{}) (out T, err error) {
	value, ok := in.(T)
	if !ok {
		err = fmt.Errorf("failed to cast item to concreate type")
	} else {
		out = value
	}
	return
}

// New will return a sqlite db connection for the databaseFile passed along.
// If the file does not exist then the an empty databasefile will be created at
// that location
// If the file does not exist and cannot be created then the an error will be
// returned
func New(ctx context.Context, databaseFile string) (db *sqlx.DB, err error) {
	slog.Debug("[datastore.New] called", slog.String("databaseFile", databaseFile))

	// if there is no error creating the database, then return the connection
	if err = createDatabaseFile(databaseFile); err == nil {
		dataSource := databaseFile + connectionParams
		db, err = sqlx.ConnectContext(ctx, driverName, dataSource)
	}
	return
}

// createDatabaseFile will look for the filepath specified, if it
// does not exist then it will create the directory path and an empty
// version of the file
// Returns an error if os.WriteFile fails, otherwise returns nil
func createDatabaseFile(databaseFile string) (err error) {

	if _, err = os.Stat(databaseFile); errors.Is(err, os.ErrNotExist) {
		// create the directory
		directory := filepath.Dir(databaseFile)
		os.MkdirAll(directory, os.ModePerm)
		// write an empty stub file to the location - if there is an error, panic
		if err = os.WriteFile(databaseFile, []byte(""), os.ModePerm); err != nil {
			slog.Error("mustCreateDatabaseFile failed", slog.String("err", err.Error()))
		}
	}
	return
}
