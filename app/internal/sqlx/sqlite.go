package sqlx

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type AccessMode string

const (
	READONLY  AccessMode = "r"
	READWRITE AccessMode = "rw"
)

const sqliteDriver string = "sqlite3" // driver name

// Sqlite used to handle creating connections.
type Sqlite struct {
	filepath string  // filepath to the database file to use for the connection
	driver   string  // driver name to use for sqlite connection
	readOnly bool    // is this readonly connection?
	db       *sql.DB // internal rw db connection

}

// Driver returns the driver name used within `sql.Open`
func (self *Sqlite) Driver() string {
	return self.driver
}

// DataSource returns the path & connection string parameters
// to use for `sql.Open`
func (self *Sqlite) DataSource() string {
	return fmt.Sprintf("%s%s", self.filepath, self.parameterString())
}

func (self *Sqlite) Mode() AccessMode {
	if self.readOnly {
		return READONLY
	}
	return READWRITE
}

// Open returns read or write DB (depending on readOnly flag) by using
// the Driver() & DataSource() functions. If a connection is already
// open, then return that.
//
// Note: will create the containing directory if it does not exist
func (self *Sqlite) Open() (db *sql.DB, err error) {
	var parentDir string = filepath.Dir(self.filepath)
	// if we've already got a connection, return that
	if self.db != nil {
		return self.db, nil
	}
	// create the parent folder
	err = os.MkdirAll(parentDir, os.ModePerm)
	if err != nil {
		return
	}
	// create the new connection
	db, err = sql.Open(self.Driver(), self.DataSource())
	if err == nil {
		self.db = db
	}
	return

}

// Close the sql connection if its present
func (self *Sqlite) Close() (err error) {
	if self.db != nil {
		err = self.db.Close()
		self.db = nil
	}

	return
}

// parameters is an internal method used to generate
// the parameters for the connection string
func (self *Sqlite) parameters() (params []string) {
	params = []string{
		"_journal=WAL",
		"_busy_timeout=5000",
		"_vacuum=incremental",
		"_synchronous=NORMAL",
		"_cache_size=1000000000",
	}
	// add on readonly
	if self.readOnly {
		params = append(params, "_readonly=1")
	}
	return params
}

// parameterString creates the parameter string section of the
// data source string (so the `?` onwards)
func (self *Sqlite) parameterString() string {
	return fmt.Sprintf("?%s",
		strings.Trim(
			strings.Join(
				self.parameters(),
				"&"),
			"&",
		),
	)

}

func NewSQLite(path string, readOnly bool) *Sqlite {
	return &Sqlite{
		filepath: path,
		driver:   sqliteDriver,
		readOnly: readOnly,
	}
}
