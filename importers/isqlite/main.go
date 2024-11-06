/*
isqlite imports data from json files into sqlite database.

Usage:

	isqlite [flags]

The flags are:

	-directory
		The source directory containing *.json files that will be
		imported.
		Each json file should be a list of objects.
	-type
		Flag to say what type of data is within these files from
		one of the known values - [costs|standards|uptime].
		Defaults to costs
	-database
		The path to the datbase file. Uses {type} as placeholder.
		Default: `./databases/{type}.db`

It will iterate over all files within the named directory, importing
each record into a database.

There are no duplication checks, it assumes this has been handled externally.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/importers/isqlite/lib"
)

var (
	args    = &lib.Arguments{}
	pattern = "*.json"
)

// Run processes all the data files in the directory and converts to a database
func Run(args *lib.Arguments) (err error) {
	var (
		files     []string
		db        *sqlx.DB
		ctx       context.Context = context.Background()
		waitgroup sync.WaitGroup  = sync.WaitGroup{}
	)
	// grab files for the directory
	if files, err = filepath.Glob(filepath.Join(args.Directory, pattern)); err != nil {
		slog.Error("[sqlite.main] file glob failed", slog.String("err", err.Error()))
		return
	}
	// get the database
	db, err = lib.GetDatabase(ctx, args)
	defer db.Close()
	if err != nil {
		slog.Error("[sqlite.main] get database failed", slog.String("err", err.Error()))
		return
	}
	// process each file in a go func and insert the content
	for _, file := range files {
		waitgroup.Add(1)
		go func(c context.Context, wg *sync.WaitGroup, d *sqlx.DB, a *lib.Arguments, f string) {
			_, e := lib.ProcessDataFile(c, d, a, f)
			if e != nil {
				err = errors.Join(err, e)
			}
			wg.Done()
		}(ctx, &waitgroup, db, args, file)
	}
	// Wait till all have been done
	waitgroup.Wait()

	return

}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[sqlite.main] init...")
	slog.Info("[sqlite.main]", slog.String("args", fmt.Sprintf("%+v", args)))

	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("arg validation failed", slog.String("err", err.Error()))
		os.Exit(1)
	}

	Run(args)
	slog.Info("[sqlite.main] done.")
}
