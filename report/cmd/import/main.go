package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domain/accounts/accountcli"
	"opg-reports/report/internal/domain/codebases/codebasecli"
	"opg-reports/report/internal/domain/codeowners/codeownercli"
	"opg-reports/report/internal/domain/infracosts/infracostcli"
	"opg-reports/report/internal/domain/teams/teamcli"
	"opg-reports/report/internal/domain/uptime/uptimecli"
	"opg-reports/report/internal/utils/logger"
	"os"

	"github.com/spf13/cobra"
)

// config items
var (
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:               "import",
	Short:             `import fetches data from api source to then populate the local database.`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// setup default values for config and logging
func init() {
	ctx = context.Background()
	log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))

}

func main() {
	var err error

	rootCmd.AddCommand(
		accountcli.CMD(ctx, log),
		codeownercli.CMD(ctx, log),
		codebasecli.CMD(ctx, log),
		infracostcli.CMD(ctx, log),
		teamcli.CMD(ctx, log),
		uptimecli.CMD(ctx, log),
	)

	err = rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
