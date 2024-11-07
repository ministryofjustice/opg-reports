package datastore

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/record"
	"github.com/ministryofjustice/opg-reports/pkg/timer"
)

// JoinOne checks that the item passed is a joinable record (record.JoinedRecord) and then calls
// that items ProcessJoins function to handle joins
// - passes along the database and transactions
func JoinOne(ctx context.Context, db *sqlx.DB, item interface{}) (err error) {
	var ok bool
	var joinable record.JoinedRecord

	// its not a joinable record, so skip
	if joinable, ok = item.(record.JoinedRecord); !ok {
		return
	}

	err = joinable.ProcessJoins(ctx, db)

	return
}

// JoinMany handles joins for multiple records at once.
// Called from within InsertMany after that transations has completed to deal with any joins
// on a model to links up data correctly from the main struct, creating records elsewhere in the
// database
// Calls JoinOne for each record
// no concurrency atm
func JoinMany[R record.Record](ctx context.Context, db *sqlx.DB, records []R) (err error) {
	var mainTimer *timer.Timer = timer.New()

	for _, record := range records {
		if e := JoinOne(ctx, db, record); e != nil {
			err = errors.Join(err, e)
		}
	}
	if err != nil {
		slog.Error("[datastore.JoinMany] error: [%s]", slog.String("err", err.Error()))
		return
	}

	mainTimer.Stop()

	return
}
