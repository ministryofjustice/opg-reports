// Package db contains interfaces for using databases and data records.
//
// Interfaces relating to connecting to the database are:
//   - Adaptor
//   - Transactional
//   - Formattable
//   - Seeder
//   - Connector
//
// Interfaces around database tables:
//   - Table
//   - Creatable
//
// Interfaces for table rows
//   - Record
//   - Row
//   - Cloneable
//   - Insertable
package dbs

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
)

// Adaptor is a complete interface for connecting to a database
type Adaptor interface {
	Connector
	Seeder
	Formattable
	Transactional
}

// Transactional interface supplies methods to provide
// a db pointer, a transaction and then commit that
// transaction
type Transactional interface {
	// GetDB returns an active pointer to a database or an error
	GetDB(ctx context.Context, driver DriverName, connection ConnectionString) (db *sqlx.DB, err error)
	// MustGetDB tries to return a db pointer, if it cant it will throw a panic
	MustGetDB(ctx context.Context, driver DriverName, connection ConnectionString) (db *sqlx.DB)
	// GetTransaction uses the data
	GetTransaction(ctx context.Context, db *sqlx.DB, readOnly bool) (tx *sqlx.Tx, err error)
	MustGetTransaction(ctx context.Context, db *sqlx.DB, readOnly bool) (tx *sqlx.Tx)

	CommitTransaction(tx *sqlx.Tx, withRollback bool) (err error)
}

// Connector has all of the methods needed to
// generate connection details
type Connector interface {
	GetDriverName() DriverName
	GetPath() DatabasePath
	GetParams() ConnectionParameters
	GetConnectionString(path DatabasePath, params ConnectionParameters) ConnectionString
}

// DriverName is a wrapper for string used to represent the name of the db driver
type DriverName string

// DatabasePath is string wrapper used to store the base connection to the database
type DatabasePath string

// ConnectionParameters is a string wrapper used to store the options for the connection string
type ConnectionParameters string

// ConnectionString is the full connection string to the database
type ConnectionString string

// Seeder contains the methods used to decide
// if the database can be seeded and then
// method to confirm it has been
type Seeder interface {
	Seedable() bool
	Seeded()
}

// Formattable exposes methods to that return
// layout / format values that can then be used in
// sql statements.
// Typically to provide the correct format to
// convert timestamp into a yyyy-mm string etc
type Formattable interface {
	DateFormat(interval dateintervals.Interval) (layout dateformats.Format)
}

// Table interface provides methods to returne details
// about the table
type Table interface {
	Table() TableName
}

// Creatable is the interface for models to create a table from the list of columns
// and then the indexes for this table as well
type Createable interface {
	// Columns returns all database columns with the column name as the
	// key and typing (INTEGER PRIMARY KEY | TEXT NOT NULL etc) as the
	// value
	Columns() map[string]string
	// Indexes should return a set of indexes to create on the table
	// with the key being the name of the index and the value
	// being a list of the fields to to use.
	Indexes() map[string][]string
}

// TableName wraps string as a return type
type TableName string

// Row interface is a baseline for a single enrtry within
// the database table (ie a row)
type Row interface {
	GetID() int
	SetID(id int)
}

// Cloneable is a special interface that exposes
// methods to generate a new, but empty version
// of itself
type Cloneable interface {
	New() Cloneable
}

// Insertable interface provides methods used to insert value into this
// table from a struct of it - provides the insert query details
type Insertable interface {
	// InsertColumns returns list of fields that should be inserted
	// and the store wrapper will work out the values
	InsertColumns() []string
}
