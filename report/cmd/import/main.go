package main

import (
	"context"
	"opg-reports/report/internal/global"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/times"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// default values for the args
var flags *global.ImportArgs

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

func getFlags(end time.Time) *global.ImportArgs {
	var start time.Time = times.ResetMonth(end)
	var startCosts time.Time = times.Add(times.ResetMonth(end), -2, times.MONTH)
	// if its first day of the month, then start should still be start of the previous month
	if end.Day() == 1 {
		start = times.Add(times.ResetMonth(end), -1, times.MONTH)
	}

	return &global.ImportArgs{
		Driver:         "sqlite3",
		DB:             "./database/api.db",
		DateEnd:        times.AsYMDString(end),
		DateStart:      times.AsYMDString(start),
		DateStartCosts: times.AsYMDString(startCosts),
		Region:         "eu-west-1",
		SrcFile:        "",
		OrgSlug:        "ministryofjustice",
		ParentSlug:     "opg",
		Filter:         "",
		Owners:         false,
		Stats:          false,
		Metrics:        false,
	}

}

func init() {
	var today = times.Today()

	flags = getFlags(today)
	root.PersistentFlags().StringVar(&flags.Driver, "driver", flags.Driver, "Database driver")
	root.PersistentFlags().StringVar(&flags.DB, "db", flags.DB, "Database path")
	root.PersistentFlags().StringVar(&flags.Params, "params", flags.Params, "Database params")
	root.PersistentFlags().StringVar(&flags.DateStart, "date-start", flags.DateStart, "Start date")
	root.PersistentFlags().StringVar(&flags.DateEnd, "date-end", flags.DateEnd, "End date")
	root.PersistentFlags().StringVar(&flags.Region, "region", flags.Region, "AWS Region")
	root.PersistentFlags().StringVar(&flags.SrcFile, "src-file", flags.SrcFile, "Source file to import data from")
	root.PersistentFlags().StringVar(&flags.OrgSlug, "org", flags.OrgSlug, "GitHub organisation")
	root.PersistentFlags().StringVar(&flags.ParentSlug, "parent", flags.ParentSlug, "GitHub parent team")

	// needs a wider range for cost stability
	root.PersistentFlags().StringVar(&flags.DateStartCosts, "date-start-costs", flags.DateStartCosts, "Start date for cost data")
	//
	root.PersistentFlags().BoolVar(&flags.Owners, "owners", flags.Owners, "fetch code owners.")
	root.PersistentFlags().BoolVar(&flags.Stats, "stats", flags.Stats, "fetch code base stats.")
	root.PersistentFlags().BoolVar(&flags.Metrics, "metrics", flags.Metrics, "fetch code base metrics (time based - heavy).")
	// extra filter
	root.PersistentFlags().StringVar(&flags.Filter, "filter", flags.Filter, "filter content - not used by everything yet.")
}
