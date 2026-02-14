package codeownercli

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codebases/codebaseselects"
	"opg-reports/report/internal/domain/codeowners/codeownergetter"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/utils/ghclients"
	"os"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	name  string = "codeowners"
	short string = `codeowners fetches and imports codeowner data from github [needs GITHUB_TOKEN].`
)

// the command flags used on the import cli tool
type cli struct {
	// data base
	DBPath   string // represents --db
	DBDriver string // represents --driver
	//
	Owner  string // represents --owner
	Parent string // represents --parent
}

// options to pass along to the getAndImport function
type options struct {
	Owner     string
	Parent    string
	Codebases []*codebasemodels.Codebase // list of codebases to fetch code ownership details about
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
	Owner:    "ministryofjustice",
	Parent:   "opg",
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
	cmd.Flags().StringVar(&flags.Owner, "owner", flags.Owner, "Owner / organisation")
	cmd.Flags().StringVar(&flags.Parent, "parent", flags.Parent, "Parent team")
	cmd.Flags().StringVar(&flags.DBPath, "db", flags.DBPath, "Database path")
	cmd.Flags().StringVar(&flags.DBDriver, "driver", flags.DBDriver, "Database driver")
}

// runCmd main runner
func runCmd(cmd *cobra.Command, args []string) (err error) {
	var repos []*codebasemodels.Codebase
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
	// fetch all codebases
	repos, err = codebaseselects.All(ctx, log, db)
	if err != nil {
		return
	}
	// aws client
	client, err = ghclients.New(ctx, log, token)
	if err != nil {
		return
	}
	return getAndImport(ctx, log, client.Repositories, db, &options{
		Parent:    flags.Parent,
		Owner:     flags.Owner,
		Codebases: repos,
	})

}

// getAndImport uses package getter to fetch and then insert data into the passed database
func getAndImport(ctx context.Context, log *slog.Logger, client codeownergetter.GitHubClient, db *sqlx.DB, params *options) (err error) {
	var (
		result []*dbstmts.Insert[*codeownermodels.Codeowner, int]
		data   []*codeownermodels.Codeowner = []*codeownermodels.Codeowner{}
		lg     *slog.Logger                 = log.With("func", "codeownercli.getAndImport", "owner", params.Owner, "parent", params.Parent)
	)

	lg.Info("starting codeowner import ...")

	// run the data getter command
	data, err = codeownergetter.GetCodeowners(ctx, log, client, &codeownergetter.Input{
		Codebases:  params.Codebases,
		ParentTeam: params.Parent,
		OrgSlug:    params.Owner,
	})
	if err != nil {
		return
	}
	// write the data
	result, err = dbsetup.Import[int](ctx, log, db, data, nil)
	if err != nil {
		return
	}
	lg.With("count", len(result)).Info("complete.")

	return
}
