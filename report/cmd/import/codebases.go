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
	"opg-reports/report/internal/domain/codeowners/codeowner"
	"opg-reports/report/internal/domain/codeowners/codeownerimports"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/utils/debugger"
	"opg-reports/report/internal/utils/ghclients"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	codebasesShortDesc string = `codebases fetches and imports code bases and code owner data from github.`
	codebasesLongDesc  string = `
codebases fetches data from github based on toek permissions, the moj org and opg team.

Truncates and then imports list of code bases and associated code owners.
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
	// db connection
	db, err = dbconnection.Connection(ctx, log, cfg.DB.Driver, cfg.DB.ConnectionString())
	if err != nil {
		return
	}
	// create client
	client, err = ghclients.New(ctx, log, cfg.Github.Token)
	if err != nil {
		log.Error("error connecting to client.", "err", err.Error())
		err = errors.Join(ErrGitHubConnFailed, err)
		return
	}
	err = codeimporter(ctx, log, client, db)

	return
}

// codeimporter calls both import roots
func codeimporter(ctx context.Context, log *slog.Logger, client *github.Client, db *sqlx.DB) (err error) {
	// close the db
	defer db.Close()

	err = dbmigrations.Migrate(ctx, log, db)
	if err != nil {
		return
	}

	res, err := codebasesImport(ctx, log, client.Teams, db)
	if err != nil {
		return
	}

	debugger.Dump(res)

	err = codeownerImport(ctx, log, client.Repositories, db, res)
	if err != nil {
		return
	}
	return
}

// codebasesImport
func codebasesImport(ctx context.Context, log *slog.Logger, client codebase.GitHubClient, db *sqlx.DB) (res []*codebasemodels.Codebase, err error) {
	var (
		result []*dbstatements.InsertStatement[*codebasemodels.Codebase, int]
		data   []*codebasemodels.Codebase = []*codebasemodels.Codebase{}
		opts   *codebase.Options          = &codebase.Options{ExcludeArchived: true}
	)
	res = []*codebasemodels.Codebase{}

	log = log.With("package", "import", "func", "codebasesImport")
	log.Info("starting codebases import command ...")

	// fetch the data
	data, err = codebase.GetCodebases(ctx, log, client, opts)
	if err != nil {
		return
	}
	// write the data
	result, err = codebaseimports.Import(ctx, log, db, data)
	if err != nil {
		return
	}
	// get the set of data that was imported
	for _, row := range result {
		res = append(res, row.Data)
	}

	log.With("count", len(result)).Info("completed.")

	return
}

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
