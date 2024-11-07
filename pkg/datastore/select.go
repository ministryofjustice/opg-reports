package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/record"
)

// Get returns a raw value from a query statments being used - this is typically a counter or the
// result of a sum operation ran against a series of rows
//
// Uses optional, ordered arguments instead of named parameter struct
func Get[R any](ctx context.Context, db *sqlx.DB, query SelectStatement, args ...interface{}) (result R, err error) {
	err = db.GetContext(ctx, &result, string(query), args...)
	return
}

// GetRecord fetches a single db result (of R record.Record) as a struct using a SelectStatement
func GetRecord[R record.Record](ctx context.Context, db *sqlx.DB, query SelectStatement, result R, args ...interface{}) (err error) {
	var row *sqlx.Row
	row = db.QueryRowxContext(ctx, string(query), args...)
	err = row.StructScan(result)
	return
}

// SelectMany runs the known statement against using the parameters as named values within them and returns the
// result as a slice of []R
// Expects multiple results - if you doing a single, use SelectOne (or Get)
func SelectMany[R record.Record](ctx context.Context, db *sqlx.DB, query NamedSelectStatement, params interface{}) (results []R, err error) {
	var statement *sqlx.NamedStmt
	// Check the parameters passed are valid for the query
	if err = ValidateParameters(params, Needs(query)); err != nil {
		slog.Error("[datastore.SelectMany] error validating parameters", slog.String("err", err.Error()), slog.String("R", fmt.Sprintf("%T", results)))
		return
	}
	if statement, err = db.PrepareNamedContext(ctx, string(query)); err == nil {
		err = statement.SelectContext(ctx, &results, params)
	} else {
		slog.Error("[datastore.SelectMany] error preparing named context", slog.String("err", err.Error()), slog.String("R", fmt.Sprintf("%T", results)))
	}

	if err != nil && err != sql.ErrNoRows {
		slog.Error("[datastore.SelectMany] error at exit", slog.String("err", err.Error()), slog.String("R", fmt.Sprintf("%T", results)))
	}
	return
}

// SelectMany runs the known statement against using the parameters as named values within them and returns the
// result as single version of R
// Should be used for single row results with NamedSelectStatement
// For multiple rows, used SelectMany
func SelectOne[R record.Record](ctx context.Context, db *sqlx.DB, query NamedSelectStatement, params interface{}) (result R, err error) {
	var statement *sqlx.NamedStmt
	var res []R = []R{}
	// Check the parameters passed are valid for the query
	if err = ValidateParameters(params, Needs(query)); err != nil {
		slog.Error("[datastore.SelectOne] error validating parameters", slog.String("err", err.Error()), slog.String("R", fmt.Sprintf("%T", result)))
		return
	}
	if statement, err = db.PrepareNamedContext(ctx, string(query)); err == nil {
		err = statement.SelectContext(ctx, &res, params)
	} else {
		slog.Error("[datastore.SelectOne] error preparing named context", slog.String("err", err.Error()), slog.String("R", fmt.Sprintf("%T", result)))
	}

	if err != nil && err != sql.ErrNoRows {
		slog.Error("[datastore.SelectOne] error at exit", slog.String("err", err.Error()), slog.String("R", fmt.Sprintf("%T", result)))
	}
	// return the first
	if len(res) > 0 {
		result = res[0]
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
