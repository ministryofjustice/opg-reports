package record

// Record interface
type Record interface {
	New() Record
	SetID(id int)
	UID() string
}
