package datastore

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/record"
	"github.com/ministryofjustice/opg-reports/pkg/timer"
)

// JoinOne checks that the item passed is a joinable record (record.JoinedRecord) and then calls
// that items ProcessJoins function to handle joins
// - passes along the database and transactions
func JoinOne(ctx context.Context, db *sqlx.DB, item interface{}, tx *sqlx.Tx) (err error) {
	var ok bool
	var joinable record.JoinedRecord

	// its not a joinable record, so skip
	if joinable, ok = item.(record.JoinedRecord); !ok {
		return
	}

	err = joinable.ProcessJoins(ctx, db, tx)

	return
}

// JoinMany handles joins for multiple records at once.
// Called from within InsertMany after that transations has completed to deal with any joins
// on a model to links up data correctly from the main struct, creating records elsewhere in the
// database
// Calls JoinOne for each record within go func for concurrency
func JoinMany[R record.Record](ctx context.Context, db *sqlx.DB, records []R) (err error) {
	var (
		mutex     *sync.Mutex    = &sync.Mutex{}
		waitgroup sync.WaitGroup = sync.WaitGroup{}
		tx        *sqlx.Tx       = db.MustBeginTx(ctx, transactionOptions)
		mainTimer *timer.Timer   = timer.New()
	)

	for _, record := range records {
		waitgroup.Add(1)
		go func(item R) {
			mutex.Lock()
			defer mutex.Unlock()
			e := JoinOne(ctx, db, item, tx)
			if e != nil {
				err = errors.Join(err, e)
			}
			waitgroup.Done()
		}(record)

	}
	waitgroup.Wait()
	if err != nil {
		slog.Error("[datastore.JoinMany] error: [%s]", slog.String("err", err.Error()))
		return
	}
	err = tx.Commit()
	mainTimer.Stop()

	if err != nil {
		tx.Rollback()
		slog.Error("[datastore.JoinMany] error commiting joins: [%s]", slog.String("err", err.Error()))
	}

	return
}
