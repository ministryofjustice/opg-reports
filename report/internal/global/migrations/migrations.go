package migrations

import (
	"context"
	"database/sql"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/conn"

	_ "github.com/mattn/go-sqlite3"
)

type Args struct {
	DB     string `json:"db"`     // --db
	Driver string `json:"driver"` // --driver
	Params string `json:"params"` // --params
}

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
	{Key: "create_codebase_stats", Stmt: create_codebase_stats},
	{Key: "create_codeowner", Stmt: create_codeowner},
	{Key: "create_codebase_metrics", Stmt: create_codebase_metrics},

	// {Key: "alter_codebase_metrics", Stmt: alter_codebase_metrics},
	{Key: "lowercase_team_name", Stmt: lowercase_team_name},
	{Key: "run_vacuum", Stmt: run_vacuum},
}

// Migrate is a wrapper around migrating all known migrations
func Migrate(ctx context.Context, flags *Args) (err error) {
	err = runMigrations(ctx, flags, migrations)
	return
}

// migrate will try to run the migrations passed along
func runMigrations(ctx context.Context, opts *Args, migrations []*Migration) (err error) {
	var (
		db  *sql.DB
		log *slog.Logger = cntxt.GetLogger(ctx).With("package", "global", "func", "runMigrations")
	)
	log.Info("starting ...", "db", opts.DB)
	// get the connection
	db, err = sql.Open(opts.Driver, conn.SqlitePath(opts.DB, opts.Params))
	if err != nil {
		return
	}
	// close at the end & write migrations
	defer db.Close()

	// now process all migrations, skipping those we've excluded from the migration file
	for _, migration := range migrations {
		// run the migration, if theres a error, fail
		if _, err = db.ExecContext(ctx, migration.Stmt); err != nil {
			log.Error("error with migration", "key", migration.Key, "err", err.Error())
			return
		}

	}
	log.Info("complete.")
	return
}
