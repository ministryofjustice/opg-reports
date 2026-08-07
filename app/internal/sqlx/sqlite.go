package sqlx

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const sqliteDriver string = "sqlite3" // driver name

// Sqlite used to handle creating connections
type Sqlite struct {
	filepath string  // filepath to the database file to use for the connection
	driver   string  // driver name to use for sqlite connection
	db       *sql.DB // internal db connection
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

// Open and return DB by using the Driver() & DataSource()
// functions. If a connection is already open, then return
// that.
//
// Note: will create the containing directory if it does not
// exist
func (self *Sqlite) Open() (db *sql.DB, err error) {
	// if we've already got a connection, return that
	if self.db != nil {
		return self.db, nil
	}
	// create the parent folder
	err = os.MkdirAll(filepath.Dir(self.filepath), os.ModePerm)
	if err != nil {
		return
	}

	// create a new connection
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
func (self *Sqlite) parameters() []string {
	return []string{
		"_journal=WAL",
		"_busy_timeout=5000",
		"_vacuum=incremental",
		"_synchronous=NORMAL",
		"_cache_size=1000000000",
	}
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

func NewSQLite(path string) *Sqlite {
	return &Sqlite{
		filepath: path,
		driver:   sqliteDriver,
	}
}
