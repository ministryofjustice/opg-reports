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
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
)

// Get the configuration data and the viper config for mapping to cli args
var conf, viperConf = config.New()

// useFixtures is flag set by cli args to decide if using fixture data
var useFixtures bool = false

// root command
var rootCmd = &cobra.Command{
	Use:   "importer",
	Short: "Importer",
	Long:  `importer imports everything.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
			ctx context.Context = context.Background()
			log *slog.Logger    = utils.Logger(conf.Log.Level, conf.Log.Type)
		)

		if useFixtures {
			log.Info("using seed data ...")
			err = seedData(ctx, log, conf)
		} else {
			log.Info("importing real data ...")
			err = realData(ctx, log, conf)
		}
		if err != nil {
			log.Error("error running imports", "error", err.Error())
			return
		}

	},
}

func seedData(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	// Seed team data
	_, err = team.Seed(ctx, log, conf, nil)
	if err != nil {
		return
	}
	// Seed awsaccounts
	_, err = awsaccount.Seed(ctx, log, conf, nil)
	if err != nil {
		return
	}
	// Seed awscosts
	_, err = awscost.Seed(ctx, log, conf, nil)
	if err != nil {
		return
	}
	return
}

// Run the importers that will use real data
func realData(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	// Import teams
	err = team.Import(ctx, log, conf)
	if err != nil {
		return
	}
	// Import accounts
	err = awsaccount.Import(ctx, log, conf)
	if err != nil {
		return
	}
	return
}

// init
// - Binds global config values to parameters
func init() {
	// bind the database.path config item
	rootCmd.PersistentFlags().StringVar(&conf.Database.Path, "database.path", conf.Database.Path, "Path to database file")
	viperConf.BindPFlag("database.path", rootCmd.PersistentFlags().Lookup("database.path"))

	// bind the github.organisation config item to the shorter --org
	rootCmd.PersistentFlags().StringVar(&conf.Github.Organisation, "github.organisation", conf.Github.Organisation, "GitHub organisation name")
	viperConf.BindPFlag("github.organisation", rootCmd.PersistentFlags().Lookup("github.organisation"))

	// bind a flag to decide if we are using fixture data or not
	rootCmd.PersistentFlags().BoolVar(&useFixtures, "fixtures", false, "When true, use a small predtermined dataset.")

}

func main() {

	rootCmd.Execute()

}
