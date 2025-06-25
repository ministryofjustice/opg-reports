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
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
)

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
		// Import teams
		err = team.Import(ctx, log, conf)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		// Import accounts
		err = awsaccount.Import(ctx, log, conf)
		if err != nil {
			slog.Error(err.Error())
			return
		}

	},
}

var conf, viperConf = config.New() // Get the configuration data and the viper config for mapping to cli args

// init
// - Binds global config values to parameters
func init() {
	// bind the database.path config item
	rootCmd.PersistentFlags().StringVar(&conf.Database.Path, "database.path", conf.Database.Path, "Path to database file")
	viperConf.BindPFlag("database.path", rootCmd.PersistentFlags().Lookup("database.path"))

	// bind the github.organisation config item to the shorter --org
	rootCmd.PersistentFlags().StringVar(&conf.Github.Organisation, "github.organisation", conf.Github.Organisation, "GitHub organisation name")
	viperConf.BindPFlag("github.organisation", rootCmd.PersistentFlags().Lookup("github.organisation"))

}

func main() {

	rootCmd.Execute()

}
