/*
importer handles adding data to the database for use by the API from real data with other commands.

Usage:

	importer [command]

Available commands:

	awscosts
	awsuptime

# Examples

`aws-vault exec <profile> -- importer awscosts --month="2025-08-01"`
*/
package main

import (
	"context"
	"log/slog"
	"os"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/awsr"
	"opg-reports/report/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// set up in the init
var (
	conf      *config.Config
	viperConf *viper.Viper
	ctx       context.Context
	log       *slog.Logger
)

// root command
var rootCmd = &cobra.Command{
	Use:               "importer",
	Short:             "Importer",
	Long:              `importer can populate database with seeded data ("seed") or new data via specific external api's.`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// awsAccountID returns the account id
func awsAccountID(client awsr.ClientSTSCaller, store awsr.RepositorySTS) (accountID string, err error) {
	caller, err := store.GetCallerIdentity(client)
	if caller != nil {
		accountID = *caller.Account
	}
	return
}

// init
func init() {
	conf, viperConf = config.New()
	ctx = context.Background()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)

	// extra options that aren't handled via config env values
	// awscosts - month to get data for

}

func main() {
	rootCmd.AddCommand(
		awscostsCmd,
		awsuptimeCmd,
		githubcodeownersCmd,
	)
	err := rootCmd.Execute()
	// fail on errir
	if err != nil {
		os.Exit(1)
	}

}
