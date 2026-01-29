package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
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
	codebasesShortDesc string = `codebases fetches and imports code bases / repositories from github.`
	codebasesLongDesc  string = `
codebases fetches and imports code bases / repositories from github bsaed on the moj org and opg team association.

Truncates and then imports list of code bases.
`
)

var (
	codebasesCmd *cobra.Command = &cobra.Command{
		Use:   "codebases",
		Short: codebasesShortDesc,
		Long:  codebasesLongDesc,
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
	_, err = codebasesImport(ctx, log, client.Teams, db)
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
