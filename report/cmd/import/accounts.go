package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/accounts/account"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"os"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	accountsShortDesc string = `accounts fetches and imports account data from opg-metadata releases.`
	accountsLongDesc  string = `
accounts fetches data from opg-metadata and imports that into the local database. Conflicts based on
the 'id' field are updated with new values.
`
)

var accountReleaseTag string = metaDataReleaseTag // release version tag to fetch account data from

var (
	accountsCmd *cobra.Command = &cobra.Command{
		Use:   "accounts",
		Short: accountsShortDesc,
		Long:  accountsLongDesc,
		RunE:  accountsRunE,
	}
)

// accountsRunE is wrapper to use with cobra
func accountsRunE(cmd *cobra.Command, args []string) (err error) {
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
	return importAccounts(ctx, log, client.Repositories, db)
}

// importAccounts imports all accounts from the opg-metadata repo released artifact
func importAccounts(ctx context.Context, log *slog.Logger, client account.GitHubClient, db *sqlx.DB) (err error) {
	var (
		result []*dbstmts.Insert[*accountmodels.Account, string]
		data   []*accountmodels.Account = []*accountmodels.Account{}
		opts   *account.Options         = &account.Options{}
		lg     *slog.Logger             = log.With("func", "import.importAccounts")
	)
	// config for the release
	opts.Tag = accountReleaseTag
	opts.DataDirectory, _ = os.MkdirTemp("", "__import-accounts-*")

	lg.Info("starting accounts import command ...")
	// fetch the data
	data, err = account.GetAwsAccountData(ctx, log, client, opts)
	if err != nil {
		return
	}
	// import
	result, err = dbsetup.Import[string](ctx, log, db, data, nil)
	if err != nil {
		return
	}
	lg.With("count", len(result)).Info("complete.")

	return
}

// add params to the command
func init() {
	accountsCmd.Flags().StringVar(&accountReleaseTag, "release-tag", accountReleaseTag, "Release to fetch account data from")

}
