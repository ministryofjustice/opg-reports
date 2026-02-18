package migrations

import (
	"context"
	"database/sql"
	"log/slog"

	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/conn"
	"opg-reports/report/package/files"
	"slices"
)

type Migration struct {
	Key  string
	Stmt string
}

type Args struct {
	DB            string `json:"db"`             // --db
	Driver        string `json:"driver"`         // --driver
	Params        string `json:"params"`         // --params
	MigrationFile string `json:"migration_file"` // --file
}

// RunMigrations will try to run the migrations passed along, skipping any that are within the migration
// json file
func Run(ctx context.Context, opts *Args, migrations []*Migration) (err error) {
	var (
		db      *sql.DB
		skipped int          = 0
		exclude []string     = []string{}
		done    []string     = []string{}
		log     *slog.Logger = cntxt.GetLogger(ctx).With("package", "global", "func", "RunMigrations")
	)
	log.Info("starting ...", "db", opts.DB, "migrations", opts.MigrationFile)
	// get the connection
	db, err = sql.Open(opts.Driver, conn.SqlitePath(opts.DB, opts.Params))
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
	for _, migration := range migrations {
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
