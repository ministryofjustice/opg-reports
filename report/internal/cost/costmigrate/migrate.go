// Package costmigrate
package costmigrate

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/files"
	"os"
	"path/filepath"
	"slices"

	_ "github.com/mattn/go-sqlite3"
)

// Options contains all the configuration required to run thie migration
// function to bring costs up to speed
type Options struct {
	DB            string `json:"db"`
	Driver        string `json:"driver"`
	Params        string `json:"params"`
	MigrationFile string `json:"migration_file"`
}

// Migrate will load the migrated file list
func Migrate(ctx context.Context, opts *Options) (err error) {
	var (
		db      *sql.DB
		skipped int          = 0
		exclude []string     = []string{}
		done    []string     = []string{}
		log     *slog.Logger = cntxt.GetLogger(ctx).WithGroup("costmigration")
	)
	log.Info("starting ...", "db", opts.DB, "migrations", opts.MigrationFile)
	// get the connection
	db, err = conn(ctx, opts)
	if err != nil {
		return
	}
	// close at the end & write migrations
	defer func() {
		db.Close()
		files.WriteAsJSON(ctx, opts.MigrationFile, done)
	}()

	// read json file if it exists, otherwise run all
	err = files.ReadJSON(ctx, opts.MigrationFile, &exclude)
	if err != nil {
		return
	}
	// now process all migrations, skipping those we've excluded from the migration file
	for _, migration := range allMigrations {
		var skip bool = slices.Contains(exclude, migration.Key)
		log.Debug("migrating ... ", "key", migration.Key, "skip?", skip)
		// if not in the excluded list, then run
		if !skip {
			// run the migration, if theres a error, fail
			if _, err = db.ExecContext(ctx, migration.Stmt); err != nil {
				log.Error("error with migration", "key", migration.Key, "err", err.Error())
				return
			}
		} else {
			skipped++
		}
		// track result for writing to file
		done = append(done, migration.Key)
	}
	log.Info("complete.", "skipped", skipped, "run", len(done))
	return
}

func conn(ctx context.Context, opts *Options) (db *sql.DB, err error) {
	// make the parent directory
	os.MkdirAll(filepath.Dir(opts.DB), os.ModePerm)
	// make the connection
	db, err = sql.Open(opts.Driver, fmt.Sprintf("%s%s", opts.DB, opts.Params))
	return
}
