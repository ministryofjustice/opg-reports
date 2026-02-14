package codebasecli

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codebases/codebasegetter"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/utils/ghclients"
	"os"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	name  string = "codebases"
	short string = `codebases fetches and imports active repositories from github.`
)

// the command flags used on the import cli tool
type cli struct {
	// data base
	DBPath   string // represents --db
	DBDriver string // represents --driver
	//
	Owner  string // represents --owner
	Parent string // represents --parent
	//
	ExcludeArchived bool
}

// options to pass along to the getAndImport function
type options struct {
	Owner           string
	Parent          string
	ExcludeArchived bool
}

// ctx / logs for the package
var (
	ctx context.Context // default context
	log *slog.Logger    // default logger
)

// default command options
var flags *cli = &cli{
	DBPath:          "./database/api.db",
	DBDriver:        "sqlite3",
	Owner:           "ministryofjustice",
	Parent:          "opg",
	ExcludeArchived: true,
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
	// aws client
	client, err = ghclients.New(ctx, log, token)
	if err != nil {
		return
	}
	return getAndImport(ctx, log, client.Teams, db, &options{
		ExcludeArchived: flags.ExcludeArchived,
		Parent:          flags.Parent,
		Owner:           flags.Owner,
	})

}

// getAndImport uses package getter to fetch and then insert data into the passed database
func getAndImport(ctx context.Context, log *slog.Logger, client codebasegetter.GitHubClient, db *sqlx.DB, params *options) (err error) {
	var (
		result []*dbstmts.Insert[*codebasemodels.Codebase, int]
		data   []*codebasemodels.Codebase = []*codebasemodels.Codebase{}
		lg     *slog.Logger               = log.With("func", "codebasecli.getAndImport", "owner", params.Owner)
	)

	lg.With("params", params).Info("starting codebase import ...")

	// run the data getter command
	data, err = codebasegetter.GetCodebases(ctx, log, client, &codebasegetter.Options{
		ExcludeArchived: params.ExcludeArchived,
		ParentTeam:      params.Parent,
		OrgSlug:         params.Owner,
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
