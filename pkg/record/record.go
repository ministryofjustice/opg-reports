package record

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Record interface
// Used for all the datastore related
type Record interface {
	New() Record
	SetID(id int)
	UID() string
}

// RecordInsertJoiner
// Used for more advanced records that have database joins
type RecordInsertJoiner interface {
	InsertJoins(ctx context.Context, db *sqlx.DB) error
}

// RecordSelectJoiner
// Used for more advanced to fetch their join information
type RecordSelectJoiner interface {
	SelectJoins(ctx context.Context, db *sqlx.DB) error
}
