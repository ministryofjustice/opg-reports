package adaptors

import (
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/fileutils"
)

// formats are a series of date formats used by sqlite to handle the various intervals
var sqliteFormats = map[dateintervals.Interval]dateformats.Format{
	dateintervals.Year:  dateformats.SqliteY,
	dateintervals.Month: dateformats.SqliteYM,
	dateintervals.Day:   dateformats.SqliteYMD,
}

const (
	// driver to use for sqlite3
	sqliteDriver string = "sqlite3"
	// standard connection string options for sqlite3
	sqliteParams string = "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000"
)

// SqliteFormatting provides methods for
// the dbs.Formatter interface that work
// for sqlite
type SqliteFormatting struct{}

// Date returns a sqlite date format that used to group / filter by time period requested.
// So 'Year' => '%Y', 'Month' => '%Y-%m'
func (self *SqliteFormatting) Date(interval dateintervals.Interval) (layout dateformats.Format) {
	var ok bool
	if layout, ok = sqliteFormats[interval]; !ok {
		layout = dateformats.SqliteYM
	}
	return
}

// Sqlite is an dbs.Adaptor that is setup for use with Sqlite3
type Sqlite struct {
	connector dbs.Connector
	mode      dbs.Moder
	seeder    dbs.Seeder
	db        dbs.DBer
	tx        dbs.Transactioner
	format    dbs.Formatter
}

func (self *Sqlite) Connector() dbs.Connector {
	return self.connector
}
func (self *Sqlite) Mode() dbs.Moder {
	return self.mode
}
func (self *Sqlite) Seed() dbs.Seeder {
	return self.seeder
}
func (self *Sqlite) DB() dbs.DBer {
	return self.db
}
func (self *Sqlite) TX() dbs.Transactioner {
	return self.tx
}
func (self *Sqlite) Format() dbs.Formatter {
	return self.format
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

// NewSqlite will provide a fresh Sqlite struct that can then be used
// for sqlx queries elsewhere.
// Will create a new database file at path if it does not exist.
func NewSqlite(path string, readOnly bool) (ds *Sqlite, err error) {
	var (
		exists    bool              = fileutils.Exists(path)
		seed      *Seed             = &Seed{seedable: !exists}
		mode      dbs.Moder         = &ReadWrite{}
		connect   *Connection       = &Connection{Path: path, Driver: sqliteDriver, Parameters: sqliteParams}
		dber      *SqlxDB           = &SqlxDB{}
		txer      *SqlxTransaction  = &SqlxTransaction{}
		formatter *SqliteFormatting = &SqliteFormatting{}
	)
	if readOnly {
		mode = &ReadOnly{}
	}

	ds = &Sqlite{
		mode:      mode,
		connector: connect,
		seeder:    seed,
		db:        dber,
		tx:        txer,
		format:    formatter,
	}

	if !exists {
		err = createDB(path)
	}

	return
}
