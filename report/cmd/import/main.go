package main

import (
	"context"
	"opg-reports/report/internal/global"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/times"
	"os"

	"github.com/spf13/cobra"
)

var today = times.Today()

// default values for the args
var flags = &global.ImportArgs{
	Driver:        "sqlite3",
	DB:            "./database/api.db",
	MigrationFile: "migrations.json",
	DateEnd:       times.AsYMDString(today),
	DateStart:     times.AsYMDString(times.Add(times.ResetMonth(today), -2, times.MONTH)),
	Region:        "eu-west-1",
	SrcFile:       "",
	OrgSlug:       "ministryofjustice",
	ParentSlug:    "opg",
}

// main root command
var root *cobra.Command = &cobra.Command{
	Use:               "import",
	Short:             `import data`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func main() {
	var err error
	var log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))
	var ctx = cntxt.AddLogger(context.Background(), log)

	root.AddCommand(
		teamsCmd,
		accountsCmd,
		costsCmd,
		uptimeCmd,
		codebasesCmd,
	)

	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}

func init() {
	root.PersistentFlags().StringVar(&flags.Driver, "driver", flags.Driver, "Database driver")
	root.PersistentFlags().StringVar(&flags.DB, "db", flags.DB, "Database path")
	root.PersistentFlags().StringVar(&flags.Params, "params", flags.Params, "Database params")
	root.PersistentFlags().StringVar(&flags.MigrationFile, "migration-file", flags.MigrationFile, "migration file")

	root.PersistentFlags().StringVar(&flags.DateStart, "date-start", flags.DateStart, "Start date")
	root.PersistentFlags().StringVar(&flags.DateEnd, "date-end", flags.DateEnd, "End date")

	root.PersistentFlags().StringVar(&flags.Region, "region", flags.Region, "AWS Region")

	root.PersistentFlags().StringVar(&flags.SrcFile, "src-file", flags.SrcFile, "Source file to import data from")

	root.PersistentFlags().StringVar(&flags.OrgSlug, "org", flags.OrgSlug, "GitHub organisation")
	root.PersistentFlags().StringVar(&flags.ParentSlug, "parent", flags.ParentSlug, "GitHub parent team")
}
