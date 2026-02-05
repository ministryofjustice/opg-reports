package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/domain/accounts/accountseeds"
	"opg-reports/report/internal/domain/codebases/codebaseseeds"
	"opg-reports/report/internal/domain/codeowners/codeownerseeds"
	"opg-reports/report/internal/domain/infracosts/infracostseeds"
	"opg-reports/report/internal/domain/teams/teamseeds"
	"opg-reports/report/internal/domain/uptime/uptimeseeds"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	seedCmdName   string = "seed" // root command name
	seedShortDesc string = `seed inserts test data into the configured database`
	seedLongDesc  string = `
seed inserts test data into the configured database; generally intended for development use only. This will also run database migrataions before inserting seed data.
`
)

var seedCmd = &cobra.Command{
	Use:   seedCmdName,
	Short: seedShortDesc,
	Long:  seedLongDesc,
	RunE:  seedRunE,
}

var seedDBPath string = "database/api.db" // represents --db

func seedRunE(cmd *cobra.Command, args []string) (err error) {
	var db *sqlx.DB
	// use the command flag value as the path
	cfg.DB.Path = seedDBPath
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
	var lg *slog.Logger = log.With("func", "db.runSeeds")

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

// setup default values for config and logging
func init() {
	seedCmd.Flags().StringVar(&seedDBPath, "db", seedDBPath, "Path to database")
}
