/*
importer imports data from json files into sqlite database.

Usage:

	importer [flags]

The flags are:

	-file
		The source json file with new data to add into the
		database.
	-type
		Flag to say what type of data is within these files from
		one of the known values.
	-database
		The path to the datbase file. Uses {type} as placeholder.
		Default: `./api.db`

It will iterate over all files within the named directory, importing
each record into a database.

There are no duplication checks, it assumes this has been handled externally.
*/
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ministryofjustice/opg-reports/importer/lib"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/models"
)

var args = &lib.Arguments{}

func Run(args *lib.Arguments) (err error) {
	var (
		adaptor dbs.Adaptor
		ctx     = context.Background()
	)

	f, ok := lib.TypeProcessors[args.Type]
	if !ok {
		err = fmt.Errorf("invalid type")
		return
	}
	adaptor, err = adaptors.NewSqlite(args.DatabasePath, false)
	defer adaptor.DB().Close()

	// remove existing dataset tracking table
	err = crud.Truncate(ctx, adaptor, &models.Dataset{})
	if err != nil {
		return
	}
	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		return
	}
	_, err = f(ctx, adaptor, args.SourceFile)

	// insert record
	isReal := []*models.Dataset{
		{Name: "real", Ts: time.Now().UTC().Format(dateformats.Full)},
	}
	if _, err = crud.Insert(ctx, adaptor, &models.Dataset{}, isReal...); err != nil {
		return
	}

	return

}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[importer] starting ...", slog.String("type", args.Type))
	slog.Debug("[importer]", slog.String("args", fmt.Sprintf("%+v", args)))

	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("[importer] arg validation failed", slog.String("err", err.Error()))
		os.Exit(1)
	}

	Run(args)
	slog.Info("[importer] done.")
}
