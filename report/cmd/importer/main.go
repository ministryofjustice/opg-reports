/*
Import data into a database

	importer
*/
package main

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/awscost"
	"github.com/ministryofjustice/opg-reports/report/internal/opgmetadata"
	"github.com/ministryofjustice/opg-reports/report/internal/s3"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
)

var (
	err             error
	conf, viperConf                 = config.New() // Get the configuration data and the viper config for mapping to cli args
	ctx             context.Context = context.Background()
	log             *slog.Logger    = utils.Logger(conf.Log.Level, conf.Log.Type)
)

var (
	s3Bucket string = "" // s3Bucket tracks the --s3-bucket arg that determines if some data is fetched from existing json or from api
)

type seedFunc func(ctx context.Context, log *slog.Logger, conf *config.Config, seeds []*sqldb.BoundStatement) (inserted []*sqldb.BoundStatement, err error)
type existingFunc func(ctx context.Context, log *slog.Logger, conf *config.Config) (err error)

// root command
var rootCmd = &cobra.Command{
	Use:               "import",
	Short:             "Import",
	Long:              `import can populate database with fixture data ("fixtures"), fetch data from pre-existing json ("existing") or new data via specific external api's.`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

var existingCmd = &cobra.Command{
	Use:   "existing",
	Short: "existing imports all known existing data files.",
	Long:  `existing imports all known data files (generally json) from a mix of sources (github, s3 buckets) that covers current and prior reporting data to ensure completeness`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// services for injection
		var (
			teamsService      = opgmetadata.Default[*team.TeamImport](ctx, log, conf)
			awsAccountService = opgmetadata.Default[*awsaccount.AwsAccountImport](ctx, log, conf)
			awsCostsService   = s3.Default[*awscost.AwsCostImport](ctx, log, conf)
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

var fixturesCmd = &cobra.Command{
	Use:   "fixtures",
	Short: "fixtures uses known data to populate the database",
	Long:  `fixtures empties and then populates the database with a series of known data sets to allow a create test instance.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		err = seedData(ctx, log, conf,
			team.Seed,
			awsaccount.Seed,
			awscost.Seed,
		)
		return
	},
}

// seedData runs seed calls
func seedData(ctx context.Context, log *slog.Logger, conf *config.Config, seeds ...seedFunc) (err error) {

	for _, lambda := range seeds {
		_, err = lambda(ctx, log, conf, nil)
		if err != nil {
			return
		}
	}

	return
}

// init
func init() {
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
