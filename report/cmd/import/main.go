package main

import (
	"context"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/env"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/times"
	"os"

	"github.com/spf13/cobra"
)

type cli struct {
	DB            string `json:"db"`             // --db
	Driver        string `json:"driver"`         // --driver
	Params        string `json:"params"`         // --params
	MigrationFile string `json:"migration_file"` // --file
	// aws
	Region string `json:"region"` // --region
	// date ranges
	DateStart string `json:"date_start"` // --start
	DateEnd   string `json:"date_end"`   // --end

}

var today = times.Today()

// default values for the args
var flags = &cli{
	Driver:        "sqlite3",
	DB:            "./database/api.db",
	MigrationFile: "migrations.json",
	Region:        "eu-west-1",
	DateEnd:       times.AsYMDString(today),
	DateStart:     times.AsYMDString(times.Add(times.ResetMonth(today), -2, times.MONTH)),
}

// main root command
var root *cobra.Command = &cobra.Command{
	Use:               "import",
	Short:             `import data`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// import command
var costsCmd = &cobra.Command{
	Use:   `costs`,
	Short: `import costs`,
	RunE:  runCostsImport,
}

func runCostsImport(cmd *cobra.Command, args []string) (err error) {
	// overwrite arg flags from env values
	if err = env.OverwriteStruct(&flags); err != nil {
		return
	}
	if err = migrate(cmd.Context(), flags); err != nil {
		return
	}
	if err = importCosts(cmd.Context(), flags); err != nil {
		return
	}

	return
}

func init() {
	root.PersistentFlags().StringVar(&flags.Driver, "driver", flags.Driver, "Database driver")
	root.PersistentFlags().StringVar(&flags.DB, "db", flags.DB, "Database path")
	root.PersistentFlags().StringVar(&flags.Params, "params", flags.Params, "Database params")
	root.PersistentFlags().StringVar(&flags.MigrationFile, "file", flags.MigrationFile, "migration file")

	root.PersistentFlags().StringVar(&flags.DateStart, "start", flags.DateStart, "Start date")
	root.PersistentFlags().StringVar(&flags.DateEnd, "end", flags.DateEnd, "End date")
	root.PersistentFlags().StringVar(&flags.Region, "region", flags.Region, "AWS Region")
}

func main() {
	var err error
	var log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))
	var ctx = cntxt.AddLogger(context.Background(), log)

	root.AddCommand(
		costsCmd,
	)

	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
