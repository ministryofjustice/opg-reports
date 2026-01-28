package main

import (
	"context"
	"log/slog"
	"opg-reports/report/conf"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/utils/logger"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

// config items
var (
	cfg       *conf.Config                // default config
	ctx       context.Context             // default context
	log       *slog.Logger                // default logger
	rootCmd   *cobra.Command              // base command
	cmdName   string          = "migrate" // root command name
	shortDesc string          = `migrate operates on the database to run all migration commands`
)
var longDesc string = `
migrate operates on the database to run all migration commands; generally intended for development use only.

environment variables that are utilised by this command:

	DB_PATH
		The file path of the database
`

func migrateFunc(cmd *cobra.Command, args []string) (err error) {
	var (
		db      *sqlx.DB
		driver  = cfg.DB.Driver
		connStr = cfg.DB.ConnectionString()
	)
	log = log.With("package", "migrate", "func", "migrateFunc")
	log.Info("starting migrate command ...")
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
	log.Info("completed.")
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
		RunE:  migrateFunc,
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
