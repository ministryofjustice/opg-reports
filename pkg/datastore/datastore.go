// Package datastore provides a consistnet database creation and wrapper to return sqlx.DB
// pointer
//
// The datastore package also provides common configurations for databases being used
// in this project
//
// # Datastore database abilitys
//
// The datastore package provides a set of statements and functions for handling
// database access, with those statments used for different funcs:
//
//   - CreateStatement: used for table and index creation; no arguments or named elements.
//   - InsertStatement: used to insert records into a table; uses named parameters (:timestamp etc)
//   - SelectStatement: used for simple select calls then return single values like a total or a count; allows ? bind vars
//   - NamedSelectStatement: used to more advanced selects that pull values from a struct; uses named params (:id etc).
//
// The built functions of datastore will required different types of statements to be passed along. This
// improves clarity on what the functions do and their use cases.
//
// Generally, setting up a databse, table and series of indexes use static sql, so make use of CreateStatement.
// Most of code that accesses the database to get information will be using NamedSelectStatements to provide
// values for WHERE, ORDER BY etc segments of the sql. These will generally come from a struct that is your
// primary database record.
//
// To use your struct with the datastore models it will need to implement the `record.Record` interface. This
// enforces a few functions to ensure the datastore methods will work correctly.
//
// If your database objects are more complex and contain joins between tables, there is a mechanism to deal with
// those joins. By implementing the `record.RecordInsertJoiner` interface calling either InsertOne or InsertMany
// will also trigger the InsertJoins function - allowing you to handle those relationships directly.
//
// Uses sqlx
package datastore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

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

var TxOptions *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}

// NewDB will return a sqlite db connection for the databaseFile passed along.
// If the file does not exist then the an empty databasefile will be created at
// that location
// If the file does not exist and cannot be created then the an error will be
// returned
func NewDB(ctx context.Context, variant *Config, databaseFile string) (db *sqlx.DB, isNew bool, err error) {
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
