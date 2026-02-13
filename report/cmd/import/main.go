package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	metaDataReleaseTag string = "v0.1.26" // used by the account & team imports as default release
	cmdName            string = "import"  // root command name
	shortDesc          string = `import fetches data from api source to then populate the local database.`
	longDesc           string = `
import fetches data from API sources and other locations such as json file artifacts to populate the local database. It
has a series of sub commands to call which will fetch data based on their domain / scope.
`
)

// errors
var (
	ErrGitHubTokenMissing = errors.New("missing github token value.")
	ErrGitHubConnFailed   = errors.New("github client failed with error.")
	ErrDBConnFailed       = errors.New("DB connection failed with error.")
)

// config items
var (
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:               cmdName,
	Short:             shortDesc,
	Long:              longDesc,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// command flags
var (
	dbPath   string = "./database/api.db" // represents --db
	dbDriver string = "sqlite3"           // represents --dbdriver
)

// dbconn helper for all sub commands to get database connection setup
func dbconn(ctx context.Context, log *slog.Logger) (db *sqlx.DB, err error) {
	// db connection
	db, err = dbconnection.Connection(ctx, log, dbDriver, dbPath)
	if err != nil {
		err = errors.Join(ErrDBConnFailed, err)
		return
	}
	err = dbsetup.Migrate(ctx, log, db)

	return
}

// ghclient helper for all sub commands to get gh client connection setup
func ghclient() (client *github.Client, err error) {
	var token = os.Getenv("GITHUB_TOKEN")
	// fail if there is no github token
	if token == "" {
		err = ErrGitHubTokenMissing
		return
	}
	// create client
	client, err = ghclients.New(ctx, log, token)
	if err != nil {
		log.Error("error connecting to client.", "err", err.Error())
		err = errors.Join(ErrGitHubConnFailed, err)
		return
	}
	return
}

// setup configures required vars for the commands
func setup() {
	ctx = context.Background()
	log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))

}

// setup default values for config and logging
func init() {
	setup()
	rootCmd.Flags().StringVar(&dbPath, "db", dbPath, "Path to database")
	rootCmd.Flags().StringVar(&dbDriver, "driver", dbDriver, "Datbase driver to use")
}

func main() {
	var err error

	rootCmd.AddCommand(
		accountsCmd,
		teamsCmd,
		codebasesCmd,
		codeownersCmd,
		infracostsCmd,
		uptimeCmd,
	)

	err = rootCmd.Execute()
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
