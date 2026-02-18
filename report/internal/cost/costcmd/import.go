package main

import (
	"opg-reports/report/internal/cost/costimport"
	"opg-reports/report/internal/cost/costmigrate"
	"opg-reports/report/package/awsclients"
	"opg-reports/report/package/awsid"
	"opg-reports/report/package/env"
	"opg-reports/report/package/times"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/spf13/cobra"
)

type importOptions struct {
	DB            string `json:"db"`             // --db
	Driver        string `json:"driver"`         // --driver
	Params        string `json:"params"`         // --params
	MigrationFile string `json:"migration_file"` // --file

	DateStart string `json:"date_start"` // --start
	DateEnd   string `json:"date_end"`   // --end
	Region    string `json:"region"`     // --region
}

var today = times.Today()

var importFlags = &importOptions{
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
	// overwrite arg flags from env values
	if e := env.OverwriteStruct(&migrationFlags); e != nil {
		return
	}
	// run the migration command
	err = costmigrate.Migrate(cmd.Context(), &costmigrate.Input{
		DB:            importFlags.DB,
		Driver:        importFlags.Driver,
		Params:        importFlags.Params,
		MigrationFile: importFlags.MigrationFile,
	})
	if err != nil {
		return
	}
	// aws client
	client, err := awsclients.New[*costexplorer.Client](cmd.Context(), importFlags.Region)
	if err != nil {
		return
	}
	// run import
	err = costimport.Import(cmd.Context(), client, &costimport.Input{
		DB:        importFlags.DB,
		Driver:    importFlags.Driver,
		Params:    importFlags.Params,
		DateStart: times.MustFromString(importFlags.DateStart),
		DateEnd:   times.MustFromString(importFlags.DateEnd),
		AccountID: awsid.AccountID(cmd.Context(), importFlags.Region),
	})
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
