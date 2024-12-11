package dbs

// CreateableTable is interface for tables that can be created in the database
type CreateableTable interface {
	Table
	Createable
}

// Table interface provides methods to returne details
// about the table
type Table interface {
	TableName() string
}

// Creatable is the interface for models to create a table from the list of columns
// and then the indexes for this table as well
type Createable interface {
	// Columns returns *ALL* database columns with the column name as the
	// key and typing (INTEGER PRIMARY KEY | TEXT NOT NULL etc) as the
	// value
	Columns() map[string]string
	// Indexes should return a set of indexes to create on the table
	// with the key being the name of the index and the value
	// being a list of the fields to to use.
	Indexes() map[string][]string
}

// Insertable interface provides methods used to insert value into this
// table from a struct of it - provides the insert query details
type Insertable interface {
	CreateableTable
	// InsertColumns returns list of fields that should be inserted
	// and the store wrapper will work out the values
	InsertColumns() []string
	// UniqueField returns a secondary unique field name on the db (like a slug)
	// or empty if the table only has primary key
	// Used for determining UPSERT / INSERT
	UniqueField() string
	// UpsertUpdate provides string for the `UPDATE SET ...` part of the upsert
	UpsertUpdate() string
}
