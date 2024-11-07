package datastore

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/record"
	"github.com/ministryofjustice/opg-reports/pkg/timer"
)

// JoinInsertOne checks that the item passed is a joinable record (record.RecordInsertJoiner) and then calls
// that items InsertJoins function to handle joins
// - passes along the database and transactions
func JoinInsertOne(ctx context.Context, db *sqlx.DB, item interface{}) (err error) {
	var ok bool
	var joinable record.RecordInsertJoiner

	// its not a joinable record, so skip
	if joinable, ok = item.(record.RecordInsertJoiner); !ok {
		return
	}

	err = joinable.InsertJoins(ctx, db)

	return
}

// JoinInsertMany handles joins for multiple records at once.
// Called from within InsertMany after that transations has completed to deal with any joins
// on a model to links up data correctly from the main struct, creating records elsewhere in the
// database
// Calls JoinInsertOne for each record
// no concurrency atm
func JoinInsertMany[R record.Record](ctx context.Context, db *sqlx.DB, records []R) (err error) {
	var mainTimer *timer.Timer = timer.New()

	for _, record := range records {
		if e := JoinInsertOne(ctx, db, record); e != nil {
			err = errors.Join(err, e)
		}
	}
	if err != nil {
		slog.Error("[datastore.JoinInsertMany] error: [%s]", slog.String("err", err.Error()))
		return
	}

	mainTimer.Stop()

	return
}
