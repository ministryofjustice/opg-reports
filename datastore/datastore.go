// Package datastore provides a consistnet database creation and wrapper to return sqlx.DB
// pointer
//
// The datastore package also provides common configurations for databases being used
// in this project
//
// Uses sqlx
package datastore

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/convert"
	"github.com/ministryofjustice/opg-reports/datastore/awscosts"
)

type Entity interface {
	*awscosts.Cost
}

// Config provides details for the databae being used that vary by driver
type Config struct {
	Connection         string // Connection is connection string settings
	DriverName         string // DriverName is the string name used in sqlx.Connect call
	YearFormat         string // YearFormat is the datetime pattern to use to return just the year (yyyy)
	YearMonthFormat    string // YearMonthFormat is the datetime pattern to return a year and month (yyyy-mm)
	YearMonthDayFormat string // YearMonthDayFormat is the datetime pattern to return a year, month and day (yyyy-mm-dd)
}

var Sqlite *Config = &Config{
	Connection:         "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000",
	DriverName:         "sqlite3",
	YearFormat:         "%Y",
	YearMonthFormat:    "%Y-%m",
	YearMonthDayFormat: "%Y-%m-%d",
}

// New will return a sqlite db connection for the databaseFile passed along.
// If the file does not exist then the an empty databasefile will be created at
// that location
// If the file does not exist and cannot be created then the an error will be
// returned
func New(ctx context.Context, variant *Config, databaseFile string) (db *sqlx.DB, isNew bool, err error) {
	slog.Debug("[datastore.New] called", slog.String("databaseFile", databaseFile))

	// if there is no error creating the database, then return the connection
	if isNew, err = createDatabaseFile(databaseFile); err == nil {
		dataSource := databaseFile + variant.Connection
		db, err = sqlx.ConnectContext(ctx, variant.DriverName, dataSource)
	}
	return
}

// createDatabaseFile will look for the filepath specified, if it
// does not exist then it will create the directory path and an empty
// version of the file
// Returns an error if os.WriteFile fails, otherwise returns nil
func createDatabaseFile(databaseFile string) (isNew bool, err error) {
	isNew = false
	if _, err = os.Stat(databaseFile); errors.Is(err, os.ErrNotExist) {
		isNew = true
		slog.Debug("[datastore.New] creating new database file", slog.String("databaseFile", databaseFile))
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

// ColumnValues finds all the values within rows passed for each of the columns, returning them
// as a map.
func ColumnValues[T Entity](rows []T, columns []string) (values map[string][]interface{}) {
	slog.Debug("[datastore.ColumnValues] called")
	values = map[string][]interface{}{}

	for _, row := range rows {
		mapped, err := convert.Map(row)
		if err != nil {
			slog.Error("to map failed", slog.String("err", err.Error()))
			continue
		}

		for _, column := range columns {
			// if not set, set it
			if _, ok := values[column]; !ok {
				values[column] = []interface{}{}
			}
			// add the value into the slice
			if rowValue, ok := mapped[column]; ok {
				// if they arent in there already
				if !slices.Contains(values[column], rowValue) {
					values[column] = append(values[column], rowValue)
				}
			}

		}
	}
	return
}
