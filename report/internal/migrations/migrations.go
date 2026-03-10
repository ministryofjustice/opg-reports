// Package migrations is used to run all database schema changes.
package migrations

import (
	"context"
	"database/sql"
	"log/slog"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/logger"

	_ "github.com/mattn/go-sqlite3"
)

type Migration struct {
	Key  string
	Stmt string
}

var migrations = []*Migration{
	{Key: "create_teams", Stmt: create_teams},
	{Key: "create_accounts", Stmt: create_accounts},
	{Key: "create_costs", Stmt: create_costs},
	{Key: "create_uptime", Stmt: create_uptime},
	{Key: "create_codebases", Stmt: create_codebases},
	// {Key: "create_codebase_stats", Stmt: create_codebase_stats},

	{Key: "lowercase_team_name", Stmt: lowercase_team_name},
	{Key: "run_vacuum", Stmt: run_vacuum},
}

// Migrate runs all the preset database migrations
func Migrate(ctx context.Context, opts *args.DB) (err error) {
	var (
		db  *sql.DB
		log *slog.Logger
	)
	ctx, log = logger.Get(ctx)
	log.Info("starting to run migrations ...", "db", opts.DB)

	// get the connection
	db, err = sql.Open(opts.Connection())
	if err != nil {
		log.Error("migrations: error connection to database", "err", err.Error())
		return
	}

	// close at the end & write migrations
	defer db.Close()
	// now process all migrations, skipping those we've excluded from the migration file
	for _, migration := range migrations {
		log.Info("running migration ... ", "migration", migration.Key)
		// run the migration, if theres a error, fail
		if _, err = db.ExecContext(ctx, migration.Stmt); err != nil {
			log.Error("error with migration", "key", migration.Key, "err", err.Error())
			return
		}

	}
	log.Info("database migrations complete.", "count", len(migrations))
	return
}
