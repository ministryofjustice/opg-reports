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

// config items
var (
	cfg     *conf.Config    // default config
	ctx     context.Context // default context
	log     *slog.Logger    // default logger
	rootCmd *cobra.Command  // base command
)

const (
	cmdName   string = "seed" // root command name
	shortDesc string = `seed inserts test data into the configured database`
	longDesc  string = `
seed inserts test data into the configured database; generally intended for development use only. This will also run database migrataions before inserting seed data.

environment variables that are utilised by this command:

	DB_PATH
		The file path of the database
`
)

func seedFunc(cmd *cobra.Command, args []string) (err error) {
	var (
		db      *sqlx.DB
		driver  string       = cfg.DB.Driver
		connStr string       = cfg.DB.ConnectionString()
		lg      *slog.Logger = log.With("func", "seed.seedFunc")
	)

	lg.Info("starting seed command ...")
	// db connection
	db, err = dbconnection.Connection(ctx, log, driver, connStr)
	if err != nil {
		return
	}
	// close the db
	defer db.Close()
	// migrate the database
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
	rootCmd = &cobra.Command{
		Use:   cmdName,
		Short: shortDesc,
		Long:  longDesc,
		RunE:  seedFunc,
	}
}

// setup default values for config and logging
func init() {
	setup()
}

func main() {
	var err error

	err = rootCmd.Execute()
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
