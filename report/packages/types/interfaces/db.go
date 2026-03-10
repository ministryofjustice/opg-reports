package interfaces

// Mappable provides way to convert struct into a map
type Mappable interface {
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
	Mappable
}

// Statement is a wrapper to get accurate sql statement
//
// Passes in the filter values to sql can be updated
// based on those if required
type Statement interface {
	// SQL uses the Filter to update and adjust the sql.
	//
	// In the method, uses convert.Between to change the
	// filter into map[string]interface{} so it can be
	// inspected for specific values
	SQL(filter Filterable) string
}

// Insertable is an alias for Mappable.
//
// Used for database record inserts, but the name
// provides clearer intent on its purpose
type Insertable Mappable

// Selectable is small interface to determine what
// funcs are required to be able to process a select
// statement and map that to a model
type Selectable interface {
	Sequence() []any
}
