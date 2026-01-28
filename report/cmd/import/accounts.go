package main

import (
	"context"
	"errors"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/accounts/account"
	"opg-reports/report/internal/domain/accounts/accountimports"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/utils/ghclients"
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

var (
	ErrAccountsTokenMissing = errors.New("missing github token value.")
	ErrAccountsConnFailed   = errors.New("github client failed with error.")
)

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
	return accountsImport(ctx, log, client.Repositories, db)
}

// accountsImport inner func called by the wrapper used by cobra
func accountsImport(ctx context.Context, log *slog.Logger, client account.GitHubClient, db *sqlx.DB) (err error) {
	var (
		result []*dbstatements.InsertStatement[*accountmodels.AwsAccount, string]
		data   []*accountmodels.AwsAccount       = []*accountmodels.AwsAccount{}
		opts   *account.GetAwsAccountDataOptions = &account.GetAwsAccountDataOptions{}
	)
	// config for the release
	opts.Tag = cfg.Accounts.Release
	opts.DataDirectory, _ = os.MkdirTemp("", "__import-accounts-*")

	log = log.With("package", "import", "func", "accountsImport")
	log.Info("starting accounts import command ...")
	// close the db
	defer db.Close()

	err = dbmigrations.Migrate(ctx, log, db)
	if err != nil {
		return
	}
	// fetch the data
	data, err = account.GetAwsAccountData(ctx, log, client, opts)
	if err != nil {
		return
	}
	// write the data
	result, err = accountimports.Import(ctx, log, db, data)
	if err != nil {
		return
	}
	log.With("count", len(result)).Info("completed.")

	return
}
