// Package dbx provides extended functions to use with sql & sqlite
//
// Adds common / consistent way for selects, inserts etc
package dbx

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Selectable is small interface to determine what
// funcs are required to be able to process a select
// statement and map that to a model
type Selectable interface {
	Sequence() []any
}

// Insertable is an alias for Mappable.
//
// Used for database record inserts, but the name
// provides clearer intent on its purpose
type Insertable interface {
	// Map should return  map of all fields of itself
	// by using json conversion
	Map() map[string]interface{}
}

// Filterable interface is used to allow a struct to
// provide values for placeholders within a sql statement.
//
// The values will be used to replace `:${name}` placeholders
// with the sql.
type Filterable interface {
	// Map should return  map of all fields of itself
	// by using json conversion
	Map() map[string]interface{}
}

// Connector provides wraooer to create a db connection with sql.Open
type Connector interface {
	Connection() (db *sql.DB)
}
