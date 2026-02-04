package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/codebases/codebase"
	"opg-reports/report/internal/domain/codebases/codebaseimports"
	"opg-reports/report/internal/domain/codebases/codebasemodels"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	codebasesShortDesc string = `codebases fetches and imports active repositories from github.`
	codebasesLongDesc  string = `
codebases fetches and imports repositories from github bsaed on the moj org and opg team association.

Truncates before importing to remove stale / inaccurate data.
`
)

var (
	codebasesOwner      string = "ministryofjustice"
	codebasesParentTeam string = "opg"
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

	err = importCodebases(ctx, log, client.Teams, db, &codebase.Options{
		ExcludeArchived: true,
		ParentTeam:      codebasesParentTeam,
		OrgSlug:         codebasesOwner,
	})
	return
}

// importCodebases imports all known, active codebases locally
func importCodebases(ctx context.Context, log *slog.Logger, client codebase.GitHubClient, db *sqlx.DB, opts *codebase.Options) (err error) {
	var (
		result []*dbstatements.InsertStatement[*codebasemodels.Codebase, int]
		data   []*codebasemodels.Codebase = []*codebasemodels.Codebase{}
		lg     *slog.Logger               = log.With("func", "import.importCodebases")
	)

	lg.Info("starting codebases import command ...")

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

	lg.With("count", len(result)).Info("complete.")

	return
}

// add params to the command
func init() {
	codebasesCmd.Flags().StringVar(&codebasesOwner, "owner", codebasesOwner, "Owner / Organisation to fetch data about.")
	codebasesCmd.Flags().StringVar(&codebasesParentTeam, "parent", codebasesParentTeam, "Limit codebases to those owned by this team and sub-teams.")

}
