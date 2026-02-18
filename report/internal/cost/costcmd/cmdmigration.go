package main

import (
	"opg-reports/report/internal/cost/costmigrations"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/package/env"

	"github.com/spf13/cobra"
)

var migrationFlags = &migrations.Args{
	Driver:        "sqlite3",
	DB:            "./database/api.db",
	MigrationFile: "migrations.json",
}

// migration command
var migrationCmd = &cobra.Command{
	Use:   `migrate`,
	Short: `run migration command`,
	RunE:  runMigration,
}

func runMigration(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&migrationFlags); e != nil {
		return
	}
	// run the migration command
	err = migrations.Run(ctx, migrationFlags, costmigrations.Migrations)
	return
}

func init() {
	migrationCmd.Flags().StringVar(&migrationFlags.Driver, "driver", migrationFlags.Driver, "Database driver")
	migrationCmd.Flags().StringVar(&migrationFlags.DB, "db", migrationFlags.DB, "Database path")
	migrationCmd.Flags().StringVar(&migrationFlags.Params, "params", migrationFlags.Params, "Database params")
	migrationCmd.Flags().StringVar(&migrationFlags.MigrationFile, "file", migrationFlags.MigrationFile, "migration file")
}
