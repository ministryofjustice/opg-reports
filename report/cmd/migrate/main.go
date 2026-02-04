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

const (
	cmdName   string = "migrate" // root command name
	shortDesc string = `migrate operates on the database to run all migration commands`
	longDesc  string = `
migrate operates on the database to run all migration commands; generally intended for development use only.
`
)

// config items
var (
	cfg *conf.Config    // default config
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   cmdName,
	Short: shortDesc,
	Long:  longDesc,
	RunE:  migrateRunE,
}

var dbPath string = "database/api.db" // represents --db

func migrateRunE(cmd *cobra.Command, args []string) (err error) {
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
	err = runMigrations(ctx, log, db)
	return
}

func runMigrations(ctx context.Context, log *slog.Logger, db *sqlx.DB) (err error) {
	var lg *slog.Logger = log.With("func", "migrate.runMigrations")

	lg.Info("starting migrate command ...")
	// migrate the database
	err = dbmigrations.Migrate(ctx, log, db)
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
