package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/migrations"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/dbx"
	"opg-reports/report/packages/logger"
	"opg-reports/report/packages/types/interfaces"
)

var ErrNoDataGetter = errors.New("no data getter function configured.")

// runner is the internal helper to wrap around import configuration.
//
// Uses generics to reduce the amount of repeated code needed within each of
// of the import packages, so calls interfaced functions in order:
//
//   - Migrates DB
//   - If present, runs getters to return the data, with each getting the previous values passed
//   - Runs filter on the last getter result
//   - Runs transform on the filtered data
//   - Runs insert
func runner[R interfaces.Insertable, T any, O any, C any](ctx context.Context, client C,
	opts *args.Import,
	sql string,
	filter interfaces.ImportFilterF[T],
	transform interfaces.ImportTransformF[R, T],
	getters ...interfaces.ImportGetterF[T, O, C]) (filtered []T, transformed []R, err error) {

	var (
		log      *slog.Logger
		previous = []O{}
		data     = []T{}
	)
	filtered = []T{}
	transformed = []R{}

	ctx, log = logger.Get(ctx)
	// call database migration
	log.Info("migrating database ...", "db", opts.DB.DB)
	err = migrations.Migrate(ctx, opts.DB)
	if err != nil {
		return
	}

	log.Info("running import ... ")
	// no data fetchers, so error
	if len(getters) == 0 {
		err = ErrNoDataGetter
		return
	}
	// run the getters and pull data
	for _, getF := range getters {
		var got interface{}
		log.Info("running data getter function ... ")
		got, err = getF(ctx, client, opts, previous...)
		if err != nil {
			return
		}
		previous = got.([]O)
		data = got.([]T)
	}

	// run the filter & transform
	log.Info("running data filtering function ... ")
	filtered = filter(ctx, data, opts.Filters)

	log.Info("running data transformation function ... ")
	transformed, err = transform(ctx, filtered, opts)
	if err != nil {
		log.Error("error with transform", "err", err.Error())
		return
	}
	// now write to db
	log.Info("running data insertion ... ")
	err = dbx.Insert(ctx, sql, transformed, opts.DB)
	if err != nil {
		log.Error("error write data during import", "err", err.Error())
		return
	}

	log.Info("import run completed.")

	return
}
