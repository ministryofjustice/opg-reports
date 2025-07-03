/*
Import data into a database

	import [COMMAND]
*/
package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/githubr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/existing"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
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

var (
	syncDB bool = false
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
		var (
			githubClient                     = githubr.DefaultClient(conf)
			s3Client                         = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
			githubStore  *githubr.Repository = githubr.Default(ctx, log, conf)
			s3Store      *awsr.Repository    = awsr.Default(ctx, log, conf)
			sqlStore     *sqlr.Repository    = sqlr.Default(ctx, log, conf)
			existService *existing.Service   = existing.Default(ctx, log, conf)
		)

		// TEAMS
		if _, err = existService.InsertTeams(githubClient.Repositories, githubStore, sqlStore); err != nil {
			return
		}
		// ACCOUNTS
		if _, err = existService.InsertAwsAccounts(githubClient.Repositories, githubStore, sqlStore); err != nil {
			return
		}
		// COSTS
		if _, err = existService.InsertAwsCosts(s3Client, s3Store, sqlStore); err != nil {
			return
		}

		return
	},
}

// seedCmd
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed inserts known test data.",
	Long:  `seed inserts known test data for use in development.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			sqlStore    *sqlr.Repository = sqlr.Default(ctx, log, conf)
			seedService *seed.Service    = seed.Default(ctx, log, conf)
		)
		// TEAMS
		if _, err = seedService.Teams(sqlStore); err != nil {
			return
		}
		// ACCOUNTS
		if _, err = seedService.AwsAccounts(sqlStore); err != nil {
			return
		}
		// COSTS
		if _, err = seedService.AwsCosts(sqlStore); err != nil {
			return
		}

		return
	},
}

// awscostsCmd imports data from the cost explorer api directyl
var awscostsCmd = &cobra.Command{
	Use:   "awscosts",
	Short: "awscosts fetches data from the cost explorer api",
	Long:  `awscosts will call the aws costexplorer api to retrieve data for period specific.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

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
	rootCmd.PersistentFlags().StringVar(&conf.Database.Path, "database.path", conf.Database.Path, "Path to local database file")
	viperConf.BindPFlag("database.path", rootCmd.PersistentFlags().Lookup("database.path"))
	// bind the github.organisation for those commands that require it
	rootCmd.PersistentFlags().StringVar(&conf.Github.Organisation, "github.organisation", conf.Github.Organisation, "GitHub organisation name")
	viperConf.BindPFlag("github.organisation", rootCmd.PersistentFlags().Lookup("github.organisation"))

	// Command specifc args
	// awscosts - sync-db
	awscostsCmd.Flags().BoolVar(&syncDB, "--sync-db", true, "When true, will download the existing database from s3 & then then upload the updated version.")
}

func main() {
	rootCmd.AddCommand(existingCmd, seedCmd, awscostsCmd)
	rootCmd.Execute()

}
