package main

import (
	costs "opg-reports/report/internal/domains/cost/importer"
	"opg-reports/report/packages/clients"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/spf13/cobra"
)

var importCosts = &cobra.Command{
	Use:   `costs`,
	Short: `import costs from aws.`,
	RunE:  importCostsF,
}
var importUptime = &cobra.Command{
	Use:   `uptime`,
	Short: `import uptime from aws.`,
	RunE:  importUptimeF,
}

// importCostsF
func importCostsF(cmd *cobra.Command, args []string) (err error) {
	var client *costexplorer.Client
	// get the client
	client, err = clients.New[*costexplorer.Client](cmd.Context(), cliFlags.Aws.Region)
	if err != nil {
		return
	}
	// set the account id
	cliFlags.Aws.AccountID = clients.AWSAccountID(cmd.Context(), cliFlags.Aws.Region)

	runner(
		cmd.Context(),
		client,
		cliFlags,
		costs.InsertStatement,
		costs.Filter,
		costs.Transform,
		costs.Get,
	)
	return
}

// importUptimeF
func importUptimeF(cmd *cobra.Command, args []string) (err error) {
	var region = "us-east-1"
	var client *costexplorer.Client
	// get the client
	client, err = clients.New[*costexplorer.Client](cmd.Context(), region)
	if err != nil {
		return
	}
	// set the account id
	cliFlags.Aws.AccountID = clients.AWSAccountID(cmd.Context(), region)

	runner(
		cmd.Context(),
		client,
		cliFlags,
		costs.InsertStatement,
		costs.Filter,
		costs.Transform,
		costs.Get,
	)
	return
}
