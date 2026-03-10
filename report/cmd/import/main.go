package main

import (
	"context"
	"log/slog"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/logger"
	"opg-reports/report/packages/times"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var cliFlags *args.Import

// main root command
var root *cobra.Command = &cobra.Command{
	Use:               "import",
	Short:             `import data`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func defaultArgs(t time.Time) (def *args.Import) {
	def = args.Default[*args.Import](t)
	return
}

func init() {
	var now = time.Now().UTC()

	cliFlags = defaultArgs(now)
	// setup all the import flags
	// - DB
	root.PersistentFlags().StringVar(&cliFlags.DB.Driver, "driver", cliFlags.DB.Driver, "database driver type.")
	root.PersistentFlags().StringVar(&cliFlags.DB.DB, "db", cliFlags.DB.DB, "database path.")
	root.PersistentFlags().StringVar(&cliFlags.DB.Params, "params", cliFlags.DB.Params, "database connection parameters.")
	// - File
	root.PersistentFlags().StringVar(&cliFlags.File.Path, "src-file", cliFlags.File.Path, "Data file location.")
	// - Filters
	root.PersistentFlags().StringVar(&cliFlags.Filters.Filter, "filter", cliFlags.Filters.Filter, "text based filter.")
	// -- Dates
	root.PersistentFlags().TimeVar(&cliFlags.Filters.Dates.Start, "date-start", cliFlags.Filters.Dates.Start, []string{times.YMD}, "start date.")
	root.PersistentFlags().TimeVar(&cliFlags.Filters.Dates.StartCosts, "date-start-costs", cliFlags.Filters.Dates.StartCosts, []string{times.YMD}, "start date for cost import.")
	root.PersistentFlags().TimeVar(&cliFlags.Filters.Dates.End, "date-end", cliFlags.Filters.Dates.End, []string{times.YMD}, "end date.")
	// - AWS
	root.PersistentFlags().StringVar(&cliFlags.Aws.Region, "region", cliFlags.Aws.Region, "AWS region.")
	// - Github
	root.PersistentFlags().StringVar(&cliFlags.Github.Organisation, "org", cliFlags.Github.Organisation, "Github organisations slug.")
	root.PersistentFlags().StringVar(&cliFlags.Github.Parent, "parent", cliFlags.Github.Parent, "Github parent team slug.")

}

func main() {
	var err error
	var log *slog.Logger
	var ctx = context.Background()
	ctx, log = logger.Get(ctx)

	root.AddCommand(
		// files
		importTeams,
		importAccounts,
		// github related
		importCodebases,
		// aws related
		importCosts,
		importUptime,
	)

	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error("error with command", "err", err.Error())
		panic("error")
		os.Exit(1)
	}

}
