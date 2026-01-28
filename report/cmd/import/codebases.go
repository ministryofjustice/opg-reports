package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/codebases/codebase"
	"opg-reports/report/internal/domain/codebases/codebaseimports"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/utils/ghclients"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	codebasesShortDesc string = `codebases fetches and imports code base data from github.`
	codebasesLongDesc  string = `
codebases fetches data from github based on toek permissions, the moj org and opg team.
`
)

var (
	codebasesCmd *cobra.Command = &cobra.Command{
		Use:   "codebases",
		Short: accountsShortDesc,
		Long:  accountsLongDesc,
		RunE:  codebasesRunE,
	}
)

// codebasesRunE is wrapper to use with cobra
func codebasesRunE(cmd *cobra.Command, args []string) (err error) {
	var client *github.Client
	var db *sqlx.DB
	// fail if there is no github token
	if cfg.Github.Token == "" {
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
	return codebasesImport(ctx, log, client.Teams, db)
}

// accountsImport inner func called by the wrapper used by cobra
func codebasesImport(ctx context.Context, log *slog.Logger, client codebase.GitHubClient, db *sqlx.DB) (err error) {
	var (
		result []*dbstatements.InsertStatement[*codebasemodels.Codebase, int]
		data   []*codebasemodels.Codebase = []*codebasemodels.Codebase{}
		opts   *codebase.Options          = &codebase.Options{ExcludeArchived: true}
	)

	log = log.With("package", "import", "func", "codebasesImport")
	log.Info("starting codebases import command ...")
	// close the db
	defer db.Close()

	err = dbmigrations.Migrate(ctx, log, db)
	if err != nil {
		return
	}
	// fetch the data
	data, err = codebase.GetCodebases(ctx, log, client, opts)
	if err != nil {
		return
	}
	// write the data
	result, err = codebaseimports.Import(ctx, log, db, data) //accountimports.Import(ctx, log, db, data)
	if err != nil {
		return
	}
	log.With("count", len(result)).Info("completed.")

	return
}
