package main

import (
	"context"
	"log/slog"
	"opg-reports/report/conf"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/domain/accounts/accountseeds"
	"opg-reports/report/internal/domain/codebases/codebaseseeds"
	"opg-reports/report/internal/domain/codeowners/codeownerseeds"
	"opg-reports/report/internal/domain/infracosts/infracostseeds"
	"opg-reports/report/internal/domain/teams/teamseeds"
	"opg-reports/report/internal/domain/uptime/uptimeseeds"
	"opg-reports/report/internal/utils/logger"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	cmdName   string = "seed" // root command name
	shortDesc string = `seed inserts test data into the configured database`
	longDesc  string = `
seed inserts test data into the configured database; generally intended for development use only. This will also run database migrataions before inserting seed data.
`
)

// config items
var (
	cfg *conf.Config    // default config
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

var rootCmd = &cobra.Command{
	Use:   cmdName,
	Short: shortDesc,
	Long:  longDesc,
	RunE:  seedRunE,
}

var dbPath string = "database/api.db" // represents --db

func seedRunE(cmd *cobra.Command, args []string) (err error) {
	var db *sqlx.DB
	// use the command flag value as the path
	cfg.DB.Path = dbPath
	// db connection
	db, err = dbconnection.Connection(ctx, log, cfg.DB.Driver, cfg.DB.ConnectionString())
	if err != nil {
		return
	}
	// close the db
	defer db.Close()
	// run the seeds
	err = runSeeds(ctx, log, db)
	if err != nil {
		return
	}

	return
}

func runSeeds(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {
	var lg *slog.Logger = log.With("func", "seed.runSeeds")

	lg.Info("starting seed command ...")
	// // migrate the database
	err = dbmigrations.Migrate(ctx, log, db)
	if err != nil {
		return
	}

	// seed teams
	_, err = teamseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed accounts
	_, err = accountseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed codebases
	_, err = codebaseseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed infracosts
	_, err = infracostseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed uptime
	_, err = uptimeseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}
	// seed codeowners
	_, err = codeownerseeds.Seed(ctx, log, db)
	if err != nil {
		return
	}

	lg.Info("complete.")
	return
}

func setup() {
	cfg = conf.New()
	ctx = context.Background()
	log = logger.New(cfg.Log.Level, cfg.Log.Type)

}

// setup default values for config and logging
func init() {
	setup()
	rootCmd.Flags().StringVar(&dbPath, "db", dbPath, "Path to database")
}

func main() {
	var err error

	err = rootCmd.Execute()
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
