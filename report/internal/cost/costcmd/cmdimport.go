package main

import (
	"opg-reports/report/internal/global/imports"
	"opg-reports/report/package/env"
	"opg-reports/report/package/times"

	"github.com/spf13/cobra"
)

var today = times.Today()

var importFlags = &imports.Args{
	Driver:        "sqlite3",
	DB:            "./database/api.db",
	MigrationFile: "migrations.json",
	Region:        "eu-west-1",
	DateEnd:       times.AsYMDString(today),
	DateStart:     times.AsYMDString(times.Add(times.ResetMonth(today), -2, times.MONTH)),
}

// import command
var importCmd = &cobra.Command{
	Use:   `import`,
	Short: `run import command`,
	RunE:  runImport,
}

func runImport(cmd *cobra.Command, args []string) (err error) {
	var ctx = cmd.Context()
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&importFlags); e != nil {
		return
	}
	err = imports.ImportCosts(ctx, importFlags)
	return
}

func init() {
	importCmd.Flags().StringVar(&importFlags.Driver, "driver", importFlags.Driver, "Database driver")
	importCmd.Flags().StringVar(&importFlags.DB, "db", importFlags.DB, "Database path")
	importCmd.Flags().StringVar(&importFlags.Params, "params", importFlags.Params, "Database params")
	importCmd.Flags().StringVar(&importFlags.MigrationFile, "file", migrationFlags.MigrationFile, "migration file")

	importCmd.Flags().StringVar(&importFlags.DateStart, "start", importFlags.DateStart, "Start date")
	importCmd.Flags().StringVar(&importFlags.DateEnd, "end", importFlags.DateEnd, "End date")
	importCmd.Flags().StringVar(&importFlags.Region, "region", importFlags.Region, "AWS Region")
}
