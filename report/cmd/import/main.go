package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/conf"
	"opg-reports/report/internal/utils/logger"
	"os"

	"github.com/spf13/cobra"
)

// config items
var (
	cfg     *conf.Config    // default config
	ctx     context.Context // default context
	log     *slog.Logger    // default logger
	rootCmd *cobra.Command  // base command
)

const (
	cmdName   string = "import" // root command name
	shortDesc string = `import fetches data from api source to then populate the local database.`
	longDesc  string = `
import fetches data from API sources and other locations such as json file artifacts to populate the local database. It
has a series of sub commands to call which will fetch data based on their domain / scope.

environment variables that are utilised by this command:

	DB_PATH
		The file path of the database
	AWS_SESSION
		The AWS active session data
	GITHUB_TOKEN
		The GitHub token sesion token
`
)

var (
	ErrGitHubTokenMissing = errors.New("missing github token value.")
	ErrGitHubConnFailed   = errors.New("github client failed with error.")
)

func setup() {
	cfg = conf.New()
	ctx = context.Background()
	log = logger.New(cfg.Log.Level, cfg.Log.Type)
	rootCmd = &cobra.Command{
		Use:               cmdName,
		Short:             shortDesc,
		Long:              longDesc,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
}

// setup default values for config and logging
func init() {
	setup()
}

func main() {
	var err error

	rootCmd.AddCommand(
		accountsCmd,
		teamsCmd,
		codebasesCmd,
		infracostsCmd,
		uptimeCmd,
	)

	err = rootCmd.Execute()
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
