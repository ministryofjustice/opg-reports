package teamcli

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/teams/teamgetter"
	"opg-reports/report/internal/domain/teams/teammodels"
	"opg-reports/report/internal/utils/ghclients"
	"os"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	name  string = "teams"
	short string = `teams fetches and imports team data from opg-metadata releases [needs GITHUB_TOKEN].`
)

// the command flags used on the import cli tool
type cli struct {
	// data base
	DBPath   string // represents --db
	DBDriver string // represents --driver
	//
	Tag string // represents --tag
}

// options to pass along to the getAndImport function
type options struct {
	Tag string
}

// ctx / logs for the package
var (
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

// default command options
var flags *cli = &cli{
	DBPath:   "./database/api.db",
	DBDriver: "sqlite3",
	Tag:      "v0.1.26",
}

// main command
var cmd = &cobra.Command{
	Use:   name,
	Short: short,
	RunE:  runCmd,
}

var ErrGitHubTokenMissing = errors.New("missing github token value.")

func CMD(c context.Context, l *slog.Logger) *cobra.Command {
	ctx = c
	log = l
	return cmd
}

func init() {
	cmd.Flags().StringVar(&flags.Tag, "tag", flags.Tag, "Release tag")
	cmd.Flags().StringVar(&flags.DBPath, "db", flags.DBPath, "Database path")
	cmd.Flags().StringVar(&flags.DBDriver, "driver", flags.DBDriver, "Database driver")
}

// runCmd main runner
func runCmd(cmd *cobra.Command, args []string) (err error) {
	var db *sqlx.DB
	var client *github.Client
	var token = os.Getenv("GITHUB_TOKEN")

	if token == "" {
		err = ErrGitHubTokenMissing
		return
	}
	// db connection
	db, err = dbconnection.Connection(ctx, log, flags.DBDriver, flags.DBPath)
	if err != nil {
		return
	}
	defer db.Close()
	// db migration before import
	err = dbsetup.Migrate(ctx, log, db)
	if err != nil {
		return
	}
	// aws client
	client, err = ghclients.New(ctx, log, token)
	if err != nil {
		return
	}
	return getAndImport(ctx, log, client.Repositories, db, &options{
		Tag: flags.Tag,
	})

}

// getAndImport uses package getter to fetch and then insert data into the passed database
func getAndImport(ctx context.Context, log *slog.Logger, client teamgetter.GitHubClient, db *sqlx.DB, params *options) (err error) {
	var (
		dir    string
		result []*dbstmts.Insert[*teammodels.Team, string]
		data   []*teammodels.Team = []*teammodels.Team{}
		lg     *slog.Logger       = log.With("func", "teamcli.getAndImport", "tag", params.Tag)
	)

	lg.With("params", params).Info("starting account import ...")
	dir, _ = os.MkdirTemp("", "__import-teams-*")
	// run the data getter command
	data, err = teamgetter.GetTeamData(ctx, log, client, &teamgetter.Options{
		Tag:           params.Tag,
		DataDirectory: dir,
	})
	if err != nil {
		return
	}
	// write the data
	result, err = dbsetup.Import[string](ctx, log, db, data, nil)
	if err != nil {
		return
	}
	lg.With("count", len(result)).Info("complete.")

	return
}
