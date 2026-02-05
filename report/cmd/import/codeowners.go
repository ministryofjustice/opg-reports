package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codebases/codebaseselects"
	"opg-reports/report/internal/domain/codeowners/codeowner"
	"opg-reports/report/internal/domain/codeowners/codeownerimports"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"

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
	codeownerOwner      string = "ministryofjustice"
	codeownerParentTeam string = "opg"
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
	var repos []*codebasemodels.Codebase
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
	// fetch codebases
	repos, err = codebaseselects.All(ctx, log, db)
	if err != nil {
		return
	}

	return importCodeowners(ctx, log, client.Repositories, db, &codeowner.Input{
		Codebases:  repos,
		ParentTeam: codeownerParentTeam,
		OrgSlug:    codeownerOwner,
	})
}

// codeownersImport inner func called by the wrapper used by cobra
func importCodeowners(ctx context.Context, log *slog.Logger, client codeowner.GitHubClient, db *sqlx.DB, opts *codeowner.Input) (err error) {
	var (
		result []*dbstmts.Insert[*codeownermodels.Codeowner, int]
		data   []*codeownermodels.Codeowner = []*codeownermodels.Codeowner{}
		lg     *slog.Logger                 = log.With("func", "import.importCodeowners")
	)
	lg.Info("starting codeowner import command ...")

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
	lg.With("count", len(result)).Info("complete.")
	return
}

// add params to the command
func init() {
	codeownersCmd.Flags().StringVar(&codeownerOwner, "owner", codeownerOwner, "Owner / Organisation to fetch data about.")
	codeownersCmd.Flags().StringVar(&codeownerParentTeam, "parent", codeownerParentTeam, "Limit codebases to those owned by this team and sub-teams.")

}
