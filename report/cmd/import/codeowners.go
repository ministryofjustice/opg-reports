package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codebases/codebaseselects"
	"opg-reports/report/internal/domain/codeowners/codeowner"
	"opg-reports/report/internal/domain/codeowners/codeownerimports"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/utils/ghclients"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	codeownersShortDesc string = `codeowners fetches and imports codeowner data from opg-metadata releases.`
	codeownersLongDesc  string = `
codeowners fetches data from opg-metadata and imports that into the local database. Conflicts based on
the 'id' field are updated with new values.
`
)

var (
	codeownersCmd *cobra.Command = &cobra.Command{
		Use:   "codeowners",
		Short: codeownersShortDesc,
		Long:  codeownersLongDesc,
		RunE:  codeownersRunE,
	}
)

// codeownersRunE is wrapper to use with cobra
func codeownersRunE(cmd *cobra.Command, args []string) (err error) {
	var client *github.Client
	var db *sqlx.DB
	var repos []*codebasemodels.Codebase
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
	repos, err = codebaseselects.All(ctx, log, db)
	if err != nil {
		return
	}

	return codeownerImport(ctx, log, client.Repositories, db, repos)
}

// codeownersImport inner func called by the wrapper used by cobra
// codeownerImport
func codeownerImport(ctx context.Context, log *slog.Logger, client codeowner.GitHubClient, db *sqlx.DB, repos []*codebasemodels.Codebase) (err error) {
	var (
		result []*dbstatements.InsertStatement[*codeownermodels.Codeowner, int]
		data   []*codeownermodels.Codeowner = []*codeownermodels.Codeowner{}
		opts   *codeowner.Input             = &codeowner.Input{Codebases: repos}
	)

	log = log.With("package", "import", "func", "codeownerImport")
	log.Info("starting codeonwer import command ...")

	// fetch the data
	data, err = codeowner.GetCodeowners(ctx, log, client, opts)
	if err != nil {
		return
	}
	// write the data
	result, err = codeownerimports.Import(ctx, log, db, data)
	if err != nil {
		return
	}
	log.With("count", len(result)).Info("completed.")
	return
}
