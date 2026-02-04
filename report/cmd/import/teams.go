package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/teams/team"
	"opg-reports/report/internal/domain/teams/teamimports"
	"opg-reports/report/internal/domain/teams/teammodels"
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

var teamReleaseTag string = metaDataReleaseTag // release version tag to fetch account data from

var (
	ErrTeamsTokenMissing = errors.New("missing github token value.")
	ErrTeamsConnFailed   = errors.New("github client failed with error.")
)

var (
	teamsCmd *cobra.Command = &cobra.Command{
		Use:   "teams",
		Short: teamsShortDesc,
		Long:  teamsLongDesc,
		RunE:  teamsRunE,
	}
)

// teamsRunE is wrapper to use with cobra
func teamsRunE(cmd *cobra.Command, args []string) (err error) {
	var client *github.Client
	var db *sqlx.DB
	// get the github client
	client, err = ghclient()
	if err != nil {
		return
	}
	// db connection
	db, err = dbconn(ctx, log)
	if err != nil {
		return
	}
	defer db.Close()

	return teamsImport(ctx, log, client.Repositories, db)
}

// teamsImport inner func called by the wrapper used by cobra
func teamsImport(ctx context.Context, log *slog.Logger, client team.GitHubClient, db *sqlx.DB) (err error) {
	var (
		result []*dbstatements.InsertStatement[*teammodels.Team, string]
		data   []*teammodels.Team = []*teammodels.Team{}
		opts   *team.Options      = &team.Options{}
		lg     *slog.Logger       = log.With("func", "import.teamsImport")
	)
	// config for the release
	opts.Tag = teamReleaseTag
	opts.DataDirectory, _ = os.MkdirTemp("", "__import-teams-*")

	lg.Info("starting teams import command ...")

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
	lg.With("count", len(result)).Info("completed.")

	return
}

// add params to the command
func init() {
	teamsCmd.Flags().StringVar(&teamReleaseTag, "release-tag", teamReleaseTag, "Release to fetch account data from")

}
