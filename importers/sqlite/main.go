/*
isqlite imports data from json files into sqlite database.

Usage:

	isqlite [flags]

The flags are:

	-file
		The source json file with new data to add into the
		database.
	-type
		Flag to say what type of data is within these files from
		one of the known values - [costs|standards|uptime|releases].
		Defaults to costs
	-database
		The path to the datbase file. Uses {type} as placeholder.
		Default: `./databases/api.db`

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
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/importers/sqlite/lib"
	"github.com/ministryofjustice/opg-reports/internal/fileutils"
)

var (
	args    = &lib.Arguments{}
	pattern = "*.json"
)

// Run processes all the data files in the directory and converts to a database
func Run(args *lib.Arguments) (err error) {
	var (
		db        *sqlx.DB
		file      string          = args.File
		ctx       context.Context = context.Background()
		waitgroup sync.WaitGroup  = sync.WaitGroup{}
	)
	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("arg validation failed", slog.String("err", err.Error()))
		return
	}

	if !fileutils.Exists(file) {
		err = fmt.Errorf("file not found: [%s]", file)
		slog.Error("[sqlite.main] file not found", slog.String("err", err.Error()))
		return
	}

	// get the database
	db, err = lib.GetDatabase(ctx, args)
	if err != nil {
		slog.Error("[sqlite.main] get database failed", slog.String("err", err.Error()))
		return
	}
	defer db.Close()
	// process each file in a go func and insert the content
	waitgroup.Add(1)
	// i dont like single letter vars.. but for an inline func
	go func(c context.Context, wg *sync.WaitGroup, d *sqlx.DB, a *lib.Arguments, f string) {
		_, e := lib.ProcessDataFile(c, d, a, f)
		if e != nil {
			err = errors.Join(err, e)
		}
		wg.Done()
	}(ctx, &waitgroup, db, args, file)

	// Wait till all have been done
	waitgroup.Wait()

	return

}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[sqlite.main] init...")
	slog.Info("[sqlite.main]", slog.String("args", fmt.Sprintf("%+v", args)))

	err = Run(args)
	if err != nil {
		panic(err)
	}
	slog.Info("[sqlite.main] done.")
}
