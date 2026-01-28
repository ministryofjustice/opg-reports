package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/teams/team"
	"opg-reports/report/internal/domain/teams/teamimports"
	"opg-reports/report/internal/domain/teams/teammodels"
	"opg-reports/report/internal/utils/ghclients"
	"os"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	teamsShortDesc string = `teams fetches and imports team data from opg-metadata releases.`
	teamsLongDesc  string = `
teams fetches data from opg-metadata and imports that into the local database. Conflicts based on
the name field are updated with new values.
`
)

var (
	ErrTeamsTokenMissing = errors.New("missing github token value.")
	ErrTeamsConnFailed   = errors.New("github client failed with error.")
)

var (
	teamsCmd *cobra.Command = &cobra.Command{
		Use:   "teams",
		Short: accountsShortDesc,
		Long:  accountsLongDesc,
		RunE:  accountsRunE,
	}
)

// teamsRunE is wrapper to use with cobra
func teamsRunE(cmd *cobra.Command, args []string) (err error) {
	var client *github.Client
	var db *sqlx.DB
	// fail if there is no github token
	if cfg.Github.Token == "" {
		log.Error("not github token found.")
		err = ErrGitHubTokenMissing
		return
	}
	// create client
	client, err = ghclients.New(ctx, log, cfg.Github.Token)
	if err != nil {
		log.Error("error connecting to client.", "err", err.Error())
		err = errors.Join(ErrGitHubConnFailed, err)
		return
	}
	// db connection
	db, err = dbconnection.Connection(ctx, log, cfg.DB.Driver, cfg.DB.ConnectionString())
	if err != nil {
		return
	}

	return teamsImport(ctx, log, client.Repositories, db)
}

// teamsImport inner func called by the wrapper used by cobra
func teamsImport(ctx context.Context, log *slog.Logger, client team.GitHubClient, db *sqlx.DB) (err error) {
	var (
		result []*dbstatements.InsertStatement[*teammodels.Team, string]
		data   []*teammodels.Team    = []*teammodels.Team{}
		opts   *team.TeamDataOptions = &team.TeamDataOptions{}
	)
	// config for the release
	opts.Tag = cfg.Teams.Release
	opts.DataDirectory, _ = os.MkdirTemp("", "__import-teams-*")

	log = log.With("package", "import", "func", "teamsImport")
	log.Info("starting teams import command ...")
	// close the db
	defer db.Close()

	err = dbmigrations.Migrate(ctx, log, db)
	if err != nil {
		return
	}
	// fetch the data
	data, err = team.GetTeamData(ctx, log, client, opts)
	if err != nil {
		return
	}
	// write the data
	result, err = teamimports.Import(ctx, log, db, data)
	if err != nil {
		return
	}
	log.With("count", len(result)).Info("completed.")

	return
}
