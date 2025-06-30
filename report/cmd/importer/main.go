/*
Import data into a database

	import [COMMAND]
*/
package main

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/service/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/service/awscost"
	"github.com/ministryofjustice/opg-reports/report/internal/service/awss3"
	"github.com/ministryofjustice/opg-reports/report/internal/service/opgmetadata"
	"github.com/ministryofjustice/opg-reports/report/internal/service/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
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
	Use:               "import",
	Short:             "Import",
	Long:              `import can populate database with fixture data ("fixtures"), fetch data from pre-existing json ("existing") or new data via specific external api's.`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// existingCmd imports all the currently know and supported previous data
// from earlier versions of reporting that are mostly stored in s3 buckets
var existingCmd = &cobra.Command{
	Use:   "existing",
	Short: "existing imports all known existing data files.",
	Long:  `existing imports all known data files (generally json) from a mix of sources (github, s3 buckets) that covers current and prior reporting data to ensure completeness`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// services for injection
		var (
			teamsService      = opgmetadata.Default[*team.TeamImport](ctx, log, conf)
			awsAccountService = opgmetadata.Default[*awsaccount.AwsAccountImport](ctx, log, conf)
			awsCostsService   = awss3.Default[*awscost.AwsCostImport](ctx, log, conf)
		)
		// TEAMS
		if err = team.Existing(ctx, log, conf, teamsService); err != nil {
			return
		}
		// AWS ACCOUNTS
		if err = awsaccount.Existing(ctx, log, conf, awsAccountService); err != nil {
			return
		}
		// AWS COSTS
		if err = awscost.Existing(ctx, log, conf, awsCostsService); err != nil {
			return
		}

		return
	},
}

// fixturesCmd creates the database with simple, known fixture data that is used for testing and dev environments
var fixturesCmd = &cobra.Command{
	Use:   "fixtures",
	Short: "fixtures uses known data to populate the database",
	Long:  `fixtures empties and then populates the database with a series of known data sets to allow a create test instance.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if _, err = team.Seed(ctx, log, conf, nil); err != nil {
			return
		}
		if _, err = awsaccount.Seed(ctx, log, conf, nil); err != nil {
			return
		}
		if _, err = awscost.Seed(ctx, log, conf, nil); err != nil {
			return
		}

		return
	},
}

// init
func init() {
	conf, viperConf = config.New()
	ctx = context.Background()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)

	// Global flags for all commands:
	// bind the database.path config item
	rootCmd.PersistentFlags().StringVar(&conf.Database.Path, "database.path", conf.Database.Path, "Path to database file")
	viperConf.BindPFlag("database.path", rootCmd.PersistentFlags().Lookup("database.path"))
	// bind the github.organisation for those commands that require it
	rootCmd.PersistentFlags().StringVar(&conf.Github.Organisation, "github.organisation", conf.Github.Organisation, "GitHub organisation name")
	viperConf.BindPFlag("github.organisation", rootCmd.PersistentFlags().Lookup("github.organisation"))

}

func main() {
	rootCmd.AddCommand(existingCmd, fixturesCmd)
	rootCmd.Execute()

}
