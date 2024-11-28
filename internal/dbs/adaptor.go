// Package db contains interfaces for using databases and data records.
package dbs

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
)

// Reader has all the base line elements needed to connect
// and read data from a database
type Adaptor interface {
	Connector() Connector
	Mode() Moder
	Seed() Seeder
	DB() DBer
	TX() Transactioner
	Format() Formatter
}

// Connector interface holds methods and details on how to connect
// to the target database
type Connector interface {
	String() string
	DriverName() string
}

// Moder stores if transactions should be in readonly
// mode or otherwise
type Moder interface {
	Read() bool
	Write() bool
}

// Seeder interface contains methods to determine if the table can be seeded
type Seeder interface {
	// Seedable returns bool determining if the table can accept seed data
	Seedable() bool
	// Seeded makrs the table as no longer seedable
	Seeded()
}

// DBer is the interface to get and close sqlx.DB struct for use in queries.
// Most queries should be using transactions
type DBer interface {
	// Get returns a db pointer based on the connection details passed
	Get(ctx context.Context, connection Connector) (db *sqlx.DB, err error)
	// Close closes the db connection and sets the pointer to nil
	Close() error
}

// Transactioner is the interface for transactions with the database
// thats used for most queries
type Transactioner interface {
	// Get returns a transaction to the database
	Get(ctx context.Context, dbp DBer, connection Connector, mode Moder) (tx *sqlx.Tx, err error)
	// Commit tries to commit the transaction
	Commit(withRollback bool) (err error)
	//
	Rollback() (err error)
}

// Formatter interface provides methods for db specific formatting within sql statements
// to handle how varying db's use differing date formats etc
type Formatter interface {
	Date(interval dateintervals.Interval) (layout string)
}
