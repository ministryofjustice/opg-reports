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

// InsertOne writes the record to the table and returns the id of that row.
//   - uses a prepared statement to run the write
//   - will return an error if either the preparation fails or if the exec errors
//   - if transaction passed is nil, a new one is created and commited
//   - if a transaction is passed, Commit is NOT executed, presumes a wrapper above is doing this
//   - if using a fresh transation, then call JoinInsertOne to deal with joins
func InsertOne[R record.Record](ctx context.Context, db *sqlx.DB, insert InsertStatement, record R, tx *sqlx.Tx) (insertedId int, err error) {
	slog.Debug("[datastore.InsertOne]")
	var (
		transaction *sqlx.Tx = tx
		stmt        string   = string(insert)
		statement   *sqlx.NamedStmt
	)
	// create own txn if we havent got one
	if tx == nil {
		transaction = db.MustBeginTx(ctx, TxOptions)
	}

	statement, err = transaction.PrepareNamedContext(ctx, stmt)
	if err != nil {
		slog.Error("[datastore.InsertOne] error preparing insert statement",
			slog.String("err", err.Error()),
			slog.String("stmt", stmt))
		return
	}

	// if the statement fails, trigger a rollback
	if err = statement.GetContext(ctx, &insertedId, record); err != nil {
		slog.Error("[datastore.InsertOne] error inserting",
			slog.String("err", err.Error()),
			slog.String("stmt", stmt))

		tx.Rollback()
		return
	}
	// if we used our own tx, then commit
	if tx == nil {
		err = transaction.Commit()
	}
	// set the ID & run joins if there is no error and we're inserting just one
	if err == nil && tx == nil {
		record.SetID(insertedId)
		JoinInsertOne(ctx, db, record)
	}

	return
}

// InsertMany utilises go func concurrency (with mutex locking) to insert mutiple entries once.
//
// Errors and insert id's are tracked and returned. An error on a particular insert does not stop the
// other inserts, but will be returned at the end.
//
// Calls JoinInsertMany to deal with any joins on the record (by checking its record.JoinInserter) and
// deals with each of those within its on transaction and loop
//
// If the commit triggers an error then a Rollback is automatically triggered
// Designed for data import steps to allow large numbers (millions) to be inserted quickly
func InsertMany[R record.Record](ctx context.Context, db *sqlx.DB, insert InsertStatement, records []R) (insertedIds []int, err error) {
	slog.Debug("[datastore.InsertMany]", slog.Int("count to insert", len(records)))

	var (
		mutex       *sync.Mutex    = &sync.Mutex{}
		waitgroup   sync.WaitGroup = sync.WaitGroup{}
		transaction *sqlx.Tx
		mainTimer   *timer.Timer = timer.New()
	)

	transaction = db.MustBeginTx(ctx, TxOptions)

	for _, record := range records {
		waitgroup.Add(1)
		// go function wrapper
		go func(item R) {
			mutex.Lock()
			defer mutex.Unlock()
			var (
				id int
				e  error
			)
			if id, e = InsertOne(ctx, db, insert, item, transaction); err != nil {
				err = errors.Join(err, e)
			} else {
				item.SetID(id)
				insertedIds = append(insertedIds, id)
			}
			waitgroup.Done()
		}(record)
	}
	// wait for all items to be done
	waitgroup.Wait()

	if err != nil {
		slog.Error("[datastore.InsertMany] error from inserts: [%s]", slog.String("err", err.Error()))
		return
	}
	err = transaction.Commit()
	mainTimer.Stop()

	if err != nil {
		transaction.Rollback()
		slog.Error("[datastore.InsertMany] error commiting inserts: [%s]", slog.String("err", err.Error()))
	}

	// if there are no errors, deal with any joins
	if err == nil {
		err = JoinInsertMany(ctx, db, records)
	}

	slog.Debug("[datastore.InsertMany] calling join handler for records")

	return
}
