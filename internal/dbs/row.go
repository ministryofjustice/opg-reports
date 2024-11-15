package dbs

// Record interface is the full version of a database row entry
type Record interface {
	Row
	Cloneable
	Insertable
}

type InsertableRow interface {
	Row
	Insertable
}

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
	Table
	// InsertColumns returns list of fields that should be inserted
	// and the store wrapper will work out the values
	InsertColumns() []string
}

type TableOfRecord interface {
	CreateableTable
	Record
}
