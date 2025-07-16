/*
Import data into a database

	import [COMMAND]
*/
package main

import (
	"context"
	"log/slog"
	"os"

	"opg-reports/report/config"
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

// optional arguments
var (
	month string = ""
)

// root command
var rootCmd = &cobra.Command{
	Use:               "import",
	Short:             "Import",
	Long:              `import can populate database with fixture data ("fixtures"), fetch data from pre-existing json ("existing") or new data via specific external api's.`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// init
func init() {
	conf, viperConf = config.New()
	ctx = context.Background()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)

	// extra options that aren't handled via config env values
	// awscosts - month to get data for
	awscostsCmd.Flags().StringVar(&month, "month", utils.LastBillingMonth(conf.Aws.BillingDate).Format(utils.DATE_FORMATS.YMD), "The month to get cost data for. (YYYY-MM-DD)")
}

func main() {
	rootCmd.AddCommand(
		existingCmd,
		seedCmd,
		awscostsCmd)
	err := rootCmd.Execute()
	// fail on errir
	if err != nil {
		os.Exit(1)
	}

}
