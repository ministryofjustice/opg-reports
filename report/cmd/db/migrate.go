package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	migrateCmdName   string = "migrate" // root command name
	migrateShortDesc string = `migrate operates on the database to run all migration commands`
	migrateLongDesc  string = `
migrate operates on the database to run all migration commands; generally intended for development use only.
`
)

var migrateCmd *cobra.Command = &cobra.Command{
	Use:   migrateCmdName,
	Short: migrateShortDesc,
	Long:  migrateLongDesc,
	RunE:  migrateRunE,
}

var migrateDBPath string = "database/api.db" // represents --db

func migrateRunE(cmd *cobra.Command, args []string) (err error) {
	var db *sqlx.DB
	// use the command flag value as the path
	cfg.DB.Path = migrateDBPath
	// db connection
	db, err = dbconnection.Connection(ctx, log, cfg.DB.Driver, cfg.DB.ConnectionString())
	if err != nil {
		return
	}
	// close the db
	defer db.Close()
	err = runMigrations(ctx, log, db)
	return
}

func runMigrations(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {
	var lg *slog.Logger = log.With("func", "db.runMigrations")

	lg.Info("starting migrate command ...")
	err = dbsetup.Migrate(ctx, log, db)
	if err != nil {
		return
	}
	lg.Info("complete.")
	return
}

// setup default values for config and logging
func init() {
	migrateCmd.Flags().StringVar(&migrateDBPath, "db", migrateDBPath, "Path to database")
}
