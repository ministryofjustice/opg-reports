// Package sqlx contains helpful additional functions to make working with the std lib sql easier.
//
// In particular, it provides a `Bind` function that handles combining data from a struct with
// a string of SQL (using field placeholders) to provide paramters for `db.QueryContext` /
// `db.ExecContext`.
//
// # Functions
//
// `Exec` is an extension of `DB.ExecContext` which uses the `Connector` interface to
// run more direct sql statements (like CREATE & DELETE) against the DB. No binding.
//
// # Interfaces
//
// A single interface (`Connector`) is provided by this package with a concreate
// implementation (Sqlite) to reduce the amount of database paramaters passed
// within functions and to ensure db connections are opened / closed consistently.
package sqlx

import "database/sql"

type Connector interface {
	// Driver provides the string name of the driver to utilise for sql.Open
	Driver() string
	// DataSource provides the connection string used by sql.Open
	DataSource() string
	// Open provides the sql.DB connection or error details.
	// Returns an existing connection if present, otherwise opens a new one
	Open() (db *sql.DB, err error)
	// Close will close the connection and remove the stored pointer so
	// the next call to Open will generate a new db
	Close() (err error)
}
