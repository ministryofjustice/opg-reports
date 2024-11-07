package datastore

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/record"
	"github.com/ministryofjustice/opg-reports/pkg/timer"
)

// JoinInsertOne checks that the item passed is a joinable record (record.JoinInserter) and then calls
// that items InsertJoins function to handle joins
// - passes along the database and transactions
func JoinInsertOne(ctx context.Context, db *sqlx.DB, item interface{}) (err error) {
	var ok bool
	var joinable record.JoinInserter

	// its not a joinable record, so skip
	if joinable, ok = item.(record.JoinInserter); !ok {
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

// JoinSelectOne checks that the item passed is a joinable record (record.JoinSelector) and then calls
// that items SelectJoins function to handle joins
// - passes along the database and transactions
func JoinSelectOne(ctx context.Context, db *sqlx.DB, item interface{}) (err error) {
	var ok bool
	var joinable record.JoinSelector

	// its not a joinable record, so skip
	if joinable, ok = item.(record.JoinSelector); !ok {
		return
	}

	err = joinable.SelectJoins(ctx, db)
	return
}

// JoinSelectMany handles joins for multiple records at once.
// Called after SelectMany to update and joins on the struct
// no concurrency atm
func JoinSelectMany[R record.Record](ctx context.Context, db *sqlx.DB, records []R) (err error) {
	var mainTimer *timer.Timer = timer.New()

	for _, record := range records {
		if e := JoinSelectOne(ctx, db, record); e != nil {
			err = errors.Join(err, e)
		}
	}
	if err != nil {
		slog.Error("[datastore.JoinSelectMany] error: [%s]", slog.String("err", err.Error()))
		return
	}

	mainTimer.Stop()

	return
}
