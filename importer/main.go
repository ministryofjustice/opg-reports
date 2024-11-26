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

	"github.com/ministryofjustice/opg-reports/importer/lib"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
)

var (
	args = &lib.Arguments{}
)

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

	_, err = f(ctx, adaptor, args.SourceFile)

	return

}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[importer] init...")
	slog.Debug("[importer]", slog.String("args", fmt.Sprintf("%+v", args)))

	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("[importer] arg validation failed", slog.String("err", err.Error()))
		os.Exit(1)
	}

	Run(args)
	slog.Info("[importer] done.")
}
