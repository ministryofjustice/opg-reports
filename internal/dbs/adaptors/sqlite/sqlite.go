// Packacge sqlite is local config for sqlite3 database usage
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/fileutils"
)

// formats are a series of date formats used by sqlite to handle the various intervals
var formats = map[dateintervals.Interval]dateformats.Format{
	dateintervals.Year:  dateformats.SqliteY,
	dateintervals.Month: dateformats.SqliteYM,
	dateintervals.Day:   dateformats.SqliteYMD,
}

const (
	driverName       string = "sqlite3"
	connectionParams string = "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000"
)

type Sqlite struct {
	db       *sqlx.DB // pointer to db connection
	driver   string   // database driver name
	path     string   // file path to this database
	params   string   // parameters to the connection string
	seedable bool     // bool to say if this was a newly created database or an existing one
	tx       *sqlx.Tx // active transation
}

// GetDriverName returns the driver name to use
func (self *Sqlite) GetDriverName() dbs.DriverName {
	return dbs.DriverName(self.driver)
}

// GetPath returns the file path used when creating this database
func (self *Sqlite) GetPath() dbs.DatabasePath {
	return dbs.DatabasePath(self.path)
}

// GetParams returns the parameters for the connection string
func (self *Sqlite) GetParams() dbs.ConnectionParameters {
	return dbs.ConnectionParameters(self.params)
}

// GetConnectionString returns string to use to connect to the database
func (self *Sqlite) GetConnectionString(path dbs.DatabasePath, params dbs.ConnectionParameters) (conn dbs.ConnectionString) {
	conn = dbs.ConnectionString(fmt.Sprintf("%s%s", path, params))
	return
}

// GetDB returns pointer to the current database.
// Will create a new connection if self.db is nil
func (self *Sqlite) GetDB(ctx context.Context, driver dbs.DriverName, connection dbs.ConnectionString) (db *sqlx.DB, err error) {
	db, err = sqlx.ConnectContext(ctx, string(driver), string(connection))
	return
}

// MustGetDB gets the db point but instead of returning an error it will
// throw a panic - not recommended to use!
func (self *Sqlite) MustGetDB(ctx context.Context, driver dbs.DriverName, connection dbs.ConnectionString) (db *sqlx.DB) {
	var err error
	db, err = self.GetDB(ctx, driver, connection)
	if err != nil {
		panic(err)
	}
	return
}

// Seedable returns a flag to say if this database is in a state
// that seed data can be inserted (ie - empty / just created)
func (self *Sqlite) Seedable() bool {
	return self.seedable
}

// Seeded marks the database as having been seeded already
func (self *Sqlite) Seeded() {
	self.seedable = false
}

// GetTransaction returns a transaction for the database passed.
// If there is already a transaction struct at the pointer, then that
// will be used instead of creating a new version.
func (self *Sqlite) GetTransaction(ctx context.Context, db *sqlx.DB, readOnly bool) (tx *sqlx.Tx, err error) {
	var options = &sql.TxOptions{
		ReadOnly:  readOnly,
		Isolation: sql.LevelDefault,
	}
	if self.tx == nil {
		self.tx, err = db.BeginTxx(ctx, options)
	}
	tx = self.tx
	return
}

// MustGetTransaction returns a transaction by calling GetTransaction,
// but if an error is returned then panic is thrown instead.
// Not recommended for production code!
func (self *Sqlite) MustGetTransaction(ctx context.Context, db *sqlx.DB, readOnly bool) (tx *sqlx.Tx) {
	var err error
	tx, err = self.GetTransaction(ctx, db, readOnly)
	if err != nil {
		panic(err)
	}
	return
}

// CommitTransaction will run the commit on the tx passed. If withRollback is true
// then if there is an error a rollback will be triggered directly
func (self *Sqlite) CommitTransaction(tx *sqlx.Tx, withRollback bool) (err error) {
	err = tx.Commit()
	if err != nil && withRollback {
		rollError := tx.Rollback()
		err = errors.Join(err, rollError)
	}
	return
}

// DateFormat returns the format to use in a sql query to represent the requested interval.
// So if you want to express a timestamp as a Year, you will get "%Y", but it assumes
// additive progression, so `Month` is actually a Year+Month ("%Y-%m")
//
// If the interval is unkonwn then a YearMonth is returned
func (self *Sqlite) DateFormat(interval dateintervals.Interval) (layout dateformats.Format) {
	var ok bool
	if layout, ok = formats[interval]; !ok {
		layout = dateformats.SqliteYM
	}
	return
}

// createDB makes a new database file at the location passed along and
// writes empty byte slice to that file so sqlite will read it correctly.
// It will also create a directory path as well.
func createDB(path string) (err error) {
	slog.Debug("[lite] creating a new database at path.", slog.String("path", path))

	directory := filepath.Dir(path)
	os.MkdirAll(directory, os.ModePerm)

	if err = os.WriteFile(path, []byte(""), os.ModePerm); err != nil {
		slog.Error("[lite] writing empty data to db file failed", slog.String("err", err.Error()))
	}
	return
}

// New will provide a fresh Sqlite struct that can then be used
// for sqlx queries elsewhere.
// Will create a new database file at path if it does not exist.
func New(path string) (ds *Sqlite, err error) {
	var exists = fileutils.Exists(path)

	ds = &Sqlite{
		path:     path,
		seedable: !exists,
		driver:   driverName,
		params:   connectionParams,
	}

	if !exists {
		err = createDB(path)
	}

	return
}
