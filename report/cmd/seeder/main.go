/*
seeder handles adding data fixture data to the database for use by the API.

Usage:

	seeder [command]

Available commands:

	all

# Examples

`seeder all`
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
	flagMonth        string = ""
	flagIncludeCosts bool   = false
)

// root command
var rootCmd = &cobra.Command{
	Use:   "seeder",
	Short: "seed inserts known test data.",
	Long: `
seed inserts known test data for use in development.

env variables used that can be adjusted:

	DATABASE_PATH
		The file path to the sqlite database that will be used
`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// init
func init() {
	conf, viperConf = config.New()
	ctx = context.Background()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)

}

func main() {
	rootCmd.AddCommand(
		allCmd,
	)
	err := rootCmd.Execute()
	// fail on errir
	if err != nil {
		os.Exit(1)
	}

}
