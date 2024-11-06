package datastore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
	"github.com/ministryofjustice/opg-reports/pkg/timer"
)

// Exec is a string used that contains a
// sql command such as CREATE TABLE or similar
// that causes a change, but returns no value
type ExecStatement string

// CreateStatement is a subtype of ExecStatement
// specifically for running create operations
type CreateStatement ExecStatement

// InsertStatement is a string used to run INSERT
// operations against the database
type InsertStatement string

// SelectStatement is a string used as enum-esque
// type contraints for sql queries that contain SELECT
// operations
type SelectStatement string

// NamedSelectStatement is a SELECT operation that
// contains :name placeholders
type NamedSelectStatement SelectStatement

// NamedParameters are structs with fields that will be converted into named values
// within statements
type NamedParameters interface{}

// Create will uses ExecContext to run the slice of createStatements passed.
// Used to create table, add indexes and so on in sequence
// Any errors in this process will trigger a panic and exit
func Create(ctx context.Context, db *sqlx.DB, create []CreateStatement) {
	slog.Debug("[datastore.Create] ")
	for _, stmt := range create {
		if _, err := db.ExecContext(ctx, string(stmt)); err != nil {
			slog.Error("error in create", slog.String("err", err.Error()))
			panic(err)
		}

	}
}

var transactionOptions *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}

// InsertOne writes the record to the table and returns the id of that row.
//   - uses a prepared statement to run the write
//   - will return an error if either the preparation fails or if the exec errors
//   - if transaction passed is nil, a new one is created and commited
//   - if a transaction is passed, Commit is NOT executed, presumes a wrapper above is doing this
func InsertOne[R any](ctx context.Context, db *sqlx.DB, insert InsertStatement, record R, tx *sqlx.Tx) (insertedId int, err error) {
	slog.Debug("[datastore.InsertOne]")
	var (
		transaction *sqlx.Tx = tx
		stmt        string   = string(insert)
		statement   *sqlx.NamedStmt
	)
	if tx == nil {
		transaction = db.MustBeginTx(ctx, transactionOptions)
	}
	statement, err = transaction.PrepareNamedContext(ctx, stmt)

	if err != nil {
		slog.Error("[datastore.InsertOne] error preparing insert statement",
			slog.String("err", err.Error()),
			slog.String("stmt", stmt))
		return
	}
	// slog.Info("[insert]", slog.String("record", fmt.Sprintf("%+v\n", record)))

	if err = statement.GetContext(ctx, &insertedId, record); err != nil {
		slog.Error("[datastore.InsertOne] error inserting",
			slog.String("err", err.Error()),
			slog.String("stmt", stmt))
		tx.Rollback()
		return
	}
	if tx == nil {
		err = transaction.Commit()
	}
	return
}

// InsertMany utilises go func concurrency (with mutex locking) to insert mutiple entries once.
//
// Errors and insert id's are tracked and returned. An error on a particular insert does not stop the
// other inserts, but will be returned at the end.
// If the commit triggers an error then a Rollback is automatically triggered
// Designed for data import steps to allow large numbers (millions) to be inserted quickly
func InsertMany[R any](ctx context.Context, db *sqlx.DB, insert InsertStatement, records []R) (insertedIds []int, err error) {
	slog.Debug("[datastore.InsertMany]", slog.Int("count to insert", len(records)))

	var (
		mutex       *sync.Mutex    = &sync.Mutex{}
		waitgroup   sync.WaitGroup = sync.WaitGroup{}
		transaction *sqlx.Tx
		mainTimer   *timer.Timer = timer.New()
	)

	transaction = db.MustBeginTx(ctx, transactionOptions)

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
	return
}

// Get returns a raw value from a query statments being used - this is typically a counter or the
// result of a sum operation ran against a series of rows
//
// Uses optional, ordered arguments instead of named parameter struct
func Get[R any](ctx context.Context, db *sqlx.DB, query SelectStatement, args ...interface{}) (result R, err error) {
	err = db.GetContext(ctx, &result, string(query), args...)
	return
}

// Select runs the known statement against using the parameters as named values within them and returns the
// result as a slice of []R
func Select[R any](ctx context.Context, db *sqlx.DB, query NamedSelectStatement, params interface{}) (results R, err error) {
	var statement *sqlx.NamedStmt
	// Check the parameters passed are valid for the query
	if err = ValidateParameters(params, Needs(query)); err != nil {
		slog.Error("[datastore.Select] error validating parameters", slog.String("err", err.Error()))
		return
	}
	if statement, err = db.PrepareNamedContext(ctx, string(query)); err == nil {
		err = statement.SelectContext(ctx, &results, params)
	} else {
		slog.Error("[datastore.Select] error preparing named context", slog.String("err", err.Error()))
	}
	if err != nil {
		slog.Error("[datastore.Select] error at exit", slog.String("err", err.Error()))
	}
	return
}

// Needs is used in part of the validate check of the named parameters and returns
// the field names the NamedSelectStatement passed in should have
// Uses a regex to find words starting with :
func Needs(query NamedSelectStatement) (needs []string) {
	var namedParamPattern string = `(?m)(:[\w-]+)`
	var prefix string = ":"
	var re = regexp.MustCompile(namedParamPattern)
	for _, match := range re.FindAllString(string(query), -1) {
		needs = append(needs, strings.TrimPrefix(match, prefix))
	}
	return
}

// ValidateParameters checks if the parameters passed meets all the required
// needs for the query being run
func ValidateParameters[P NamedParameters](params P, needs []string) (err error) {
	mapped, err := convert.Map(params)
	if err != nil {
		return
	}
	if len(mapped) == 0 {
		err = fmt.Errorf("parameters passed must contain at least one field")
		return
	}

	missing := []string{}
	// check each need if that exists as a key in the map
	for _, need := range needs {
		if _, ok := mapped[need]; !ok {
			missing = append(missing, need)
		}
	}
	// if any field is missing then set error
	if len(missing) > 0 {
		cols := strings.Join(missing, ",")
		err = fmt.Errorf("missing required fields for this query: [%s]", cols)
	}

	return
}

// ColumnValues finds all the unique values within rows passed for each of the columns, returning them
// as a map.
func ColumnValues[T any](rows []T, columns []string) (values map[string][]interface{}) {
	slog.Debug("[datastore.ColumnValues] called")
	values = map[string][]interface{}{}

	for _, row := range rows {
		mapped, err := convert.Map(row)
		if err != nil {
			slog.Error("[datastore.ColumnValues] to map failed", slog.String("err", err.Error()))
			return
		}

		for _, column := range columns {
			// if not set, set it
			if _, ok := values[column]; !ok {
				values[column] = []interface{}{}
			}
			// add the value into the slice
			if rowValue, ok := mapped[column]; ok {
				// if they arent in there already
				if !slices.Contains(values[column], rowValue) {
					values[column] = append(values[column], rowValue)
				}
			}

		}
	}
	return
}

// Setup will ensure a database with records exists in the filepath requested.
// If there is no database at that location a new sqlite database will
// be created and populated with series of dummy data - helpful for local testing.
func Setup[T any](ctx context.Context, dbFilepath string, insertStmt InsertStatement, creates []CreateStatement, seed bool, n int) {

	var err error
	var db *sqlx.DB
	var isNew bool = false
	// add custom fakers
	exfaker.AddProviders()

	db, isNew, err = CreateNewDB(ctx, dbFilepath, creates)
	defer db.Close()

	if err != nil {
		panic(err)
	}

	if seed && isNew {
		faked := exfaker.Many[T](n)
		_, err = InsertMany(ctx, db, insertStmt, faked)
	}
	if err != nil {
		panic(err)
	}

}

// CreateNewDB will create a new DB file and then
// try to run table and index creates against it.
// Returns the db, a bool to say if it was new and any errors
func CreateNewDB(ctx context.Context, dbFilepath string, creates []CreateStatement) (db *sqlx.DB, isNew bool, err error) {

	db, isNew, err = NewDB(ctx, Sqlite, dbFilepath)
	if err == nil && len(creates) > 0 {
		Create(ctx, db, creates)
	}

	return
}
