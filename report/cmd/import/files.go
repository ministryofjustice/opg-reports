package main

import (
	accounts "opg-reports/report/internal/domains/account/importer"
	teams "opg-reports/report/internal/domains/team/importer"

	"github.com/spf13/cobra"
)

var importTeams = &cobra.Command{
	Use:   `teams`,
	Short: `import teams from local file.`,
	RunE:  importTeamsF,
}

var importAccounts = &cobra.Command{
	Use:   `accounts`,
	Short: `import accounts from local file.`,
	RunE:  importAccountsF,
}

// importTeamsF
func importTeamsF(cmd *cobra.Command, args []string) (err error) {

	runner(
		cmd.Context(),
		nil,
		cliFlags,
		teams.InsertStatement,
		teams.Filter,
		teams.Transform,
		teams.Get,
	)
	return
}

// importAccountsF
func importAccountsF(cmd *cobra.Command, args []string) (err error) {

	runner(
		cmd.Context(),
		nil,
		cliFlags,
		accounts.InsertStatement,
		accounts.Filter,
		accounts.Transform,
		accounts.Get,
	)
	return
}
