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

// JoinedRecord
// Used for more advanced records that have database joins
type JoinedRecord interface {
	ProcessJoins(ctx context.Context, db *sqlx.DB, tx *sqlx.Tx) error
}
