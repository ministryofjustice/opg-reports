package awscosts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/timer"
)

const insertCosts string = `INSERT INTO aws_costs(
    ts,
    organisation,
    account_id,
    account_name,
    unit,
    label,
    environment,
    service,
    region,
    date,
    cost
) VALUES (
    :ts,
	:organisation,
	:account_id,
	:account_name,
	:unit,
	:label,
	:environment,
	:service,
	:region,
	:date,
	:cost
) RETURNING id`

// InsertOne writes the values of cost into the database and returns the id of that row
// It uses a prepared statement to run the write and will return an error if either
// the preparation failes or if the exec errors
func InsertOne(ctx context.Context, db *sqlx.DB, cost *Cost) (insertedId int, err error) {
	slog.Debug("[awscosts.InsertOne]")
	var statement *sqlx.NamedStmt

	statement, err = db.PrepareNamedContext(ctx, insertCosts)
	if err != nil {
		slog.Error("[awscosts.InsetOne] error preparing awscosts insert statment", slog.String("err", err.Error()))
		return
	}
	err = statement.GetContext(ctx, &insertedId, cost)

	return
}

// InsertMany utilises go func concurrency (with mutex locking) to generate a series of transations
// to insert mutiple entries at the same time in a more performant fashion then looping and calling
// Insert.
// Errors and insert id's are tracked and returned. An error on a particular insert does not stop the
// other inserts, but will be returned at the end.
// If the number of inserted id's does not match the length of the items passed, this will generate
// an error as well
// If any error is found then a Rollback is automatically triggered
//
// Designed for data import steps to allow large numbers (millions) to be inserted quickly (< 60 seconds)
func InsertMany(ctx context.Context, db *sqlx.DB, costs []*Cost) (insertedIds []int, err error) {
	slog.Debug("[awscosts.InsertMany]", slog.Int("count to insert", len(costs)))
	var errs []error = []error{}

	var transactionOptions *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}
	var transaction *sqlx.Tx
	var statement *sqlx.NamedStmt

	var mutex = &sync.Mutex{}
	var waitgroup = sync.WaitGroup{}

	var mainTimer *timer.Timer = timer.New()

	transaction = db.MustBeginTx(ctx, transactionOptions)

	statement, err = transaction.PrepareNamedContext(ctx, insertCosts)
	if err != nil {
		slog.Error("error preparing awscosts insert statment for multiples", slog.String("err", err.Error()))
		return
	}

	for _, cost := range costs {
		// add to waitgroup
		waitgroup.Add(1)
		// in the function we need to lock the resource
		go func(item *Cost) {
			mutex.Lock()
			var id int
			var tick *timer.Timer = timer.New()
			// run the statement and add the inserted id into the stack
			if err = statement.GetContext(ctx, &id, item); err == nil {
				insertedIds = append(insertedIds, id)
			} else {
				errs = append(errs, err)
			}
			tick.Stop()
			mutex.Unlock()
			waitgroup.Done()
		}(cost)
	}

	// wait for all items to be done
	waitgroup.Wait()

	// commit the transaction
	err = transaction.Commit()
	if err != nil {
		slog.Error("[awscosts.InsertMany] error commiting inserts: [%s]", slog.String("err", err.Error()))
	}

	// stop the timer and output the duration details
	mainTimer.Stop()
	slog.Info("[awscosts.InsertMany] timer", slog.Int("count", len(costs)), slog.Float64("duration", mainTimer.Duration()))

	// an error has happened in the go func loops, merge all of them into error
	if len(errs) > 0 {
		merged := errors.Join(errs...)
		err = errors.Join(err, merged)
	}
	// if the count between the requested and actual items dont match, set an error
	if len(insertedIds) != len(costs) {
		mismatch := fmt.Errorf("[awscosts.InsertMany] error inserting multiple records - expected [%d] inserts, buy only have [%v] ids", len(costs), len(insertedIds))
		err = errors.Join(err, mismatch)
	}

	// rollback if any error if found
	if err != nil {
		slog.Info("[awscosts.InsertMany] errors found in multiple inserts, rolling back")
		transaction.Rollback()
	}

	return
}
