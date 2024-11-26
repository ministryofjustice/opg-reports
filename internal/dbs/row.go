package dbs

// Record interface is the full version of a database row entry
type Record interface {
	Row
	Cloneable
	Insertable
}

// InsertableRow
type InsertableRow interface {
	Row
	Insertable
}

// Row interface is a baseline for a single enrtry within
// the database table (ie a row)
type Row interface {
	GetID() int
	SetID(id int)
	//
	UniqueValue() string
}

// Cloneable is a special interface that exposes
// methods to generate a new, but empty version
// of itself
type Cloneable interface {
	New() Cloneable
}

// TableOfRecord is interface capturing both table and row
// setup needs
type TableOfRecord interface {
	CreateableTable
	Record
}
