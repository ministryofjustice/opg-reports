package datastore

import (
	"context"
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

// List returns muliple rows with a standard select - something like a select *
func List[R any](ctx context.Context, db *sqlx.DB, query SelectStatement, r R, args ...interface{}) (result []R, err error) {
	var rows *sqlx.Rows
	result = []R{}
	rows, err = db.QueryxContext(ctx, string(query), args...)
	if err != nil {
		slog.Error("[datastore.List] error calling queryx", slog.String("err", err.Error()))
		return
	}
	for rows.Next() {
		err = rows.StructScan(r)
		if err != nil {
			slog.Error("[datastore.List] error scanning row", slog.String("err", err.Error()))
			return
		}
		result = append(result, r)
	}

	return
}

// Select runs the known statement against using the parameters as named values within them and returns the
// result as a slice of []R
// Expects multiple results - if you doing a single, use Get
func Select[R record.Record](ctx context.Context, db *sqlx.DB, query NamedSelectStatement, params interface{}) (results []R, err error) {
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
