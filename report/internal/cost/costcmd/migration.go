package main

import (
	"opg-reports/report/internal/cost/costmigrate"
	"opg-reports/report/package/env"

	"github.com/spf13/cobra"
)

type migrationOptions struct {
	DB            string `json:"db"`             // --db
	Driver        string `json:"driver"`         // --driver
	Params        string `json:"params"`         // --params
	MigrationFile string `json:"migration_file"` // --file
}

var migrationFlags = &migrationOptions{
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
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&migrationFlags); e != nil {
		return
	}
	// run the migration command
	err = costmigrate.Migrate(cmd.Context(), &costmigrate.Input{
		DB:            migrationFlags.DB,
		Driver:        migrationFlags.Driver,
		Params:        migrationFlags.Params,
		MigrationFile: migrationFlags.MigrationFile,
	})
	return
}

func init() {
	migrationCmd.Flags().StringVar(&migrationFlags.Driver, "driver", migrationFlags.Driver, "Database driver")
	migrationCmd.Flags().StringVar(&migrationFlags.DB, "db", migrationFlags.DB, "Database path")
	migrationCmd.Flags().StringVar(&migrationFlags.Params, "params", migrationFlags.Params, "Database params")
	migrationCmd.Flags().StringVar(&migrationFlags.MigrationFile, "file", migrationFlags.MigrationFile, "migration file")
}
