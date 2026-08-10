// Package sqlx contains helpful additional functions to make working with the std lib sql easier.
//
// # Functions
//
// `Exec` is an extension of `DB.ExecContext` which uses the `Connector` interface to
// run more direct sql statements (like CREATE & DELETE) against the DB. No binding.
//
// `Insert` ...
//
// `Select` ...
//
// `Bind` ...
package sqlx

import (
	"database/sql"
	"errors"
	"reflect"
)

var (
	ErrReadOnlyMode     error = errors.New("connection set as readonly, requires read-write.")                                                                 // raised if the connection is readonly when write is required (exec / insert).
	ErrModelNotStruct   error = errors.New("model passed to Bind is not a struct.")                                                                            // raised when a non-struct is passed to `Bind`
	ErrModelMissingTags error = errors.New("model is missing json tags on fields.")                                                                            // raised when the struct passed to `Bind` has fields without json tags
	ErrBindingNoKey     error = errors.New("error when binding sql statement to model, a parameter (`:x`) was found that has no coresponding value on model.") // raised when the sql passed to `Bind` contains a placeholder that does not exist on the model.
)

type Connector interface {
	// Driver provides the string name of the driver to utilise for sql.Open
	Driver() string
	// DataSource provides the connection string used by sql.Open
	DataSource() string
	// Mode returns if this is readonly or readwrite
	Mode() AccessMode
	// Open provides the connection or error details.
	// Returns an existing connection if present, otherwise opens a new one
	Open() (db *sql.DB, err error)
	// Close will close the connection and remove the stored pointer so
	// the next call to Open will generate a new db
	Close() (err error)
}

// instanceOf generate new instance of T - presuming a pointer
func instanceOf[T Result]() (instance T) {
	var (
		x T
		t reflect.Type = reflect.TypeOf(x)
	)
	instance = x
	if t != nil && t.Kind() == reflect.Ptr {
		instance = reflect.New(t.Elem()).Interface().(T)
	}
	return
}
