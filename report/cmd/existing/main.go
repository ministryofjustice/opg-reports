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
	flagMonth        string = ""
	flagIncludeCosts bool   = false
)

// root command
var rootCmd = &cobra.Command{
	Use:               "existing",
	Short:             "existing",
	Long:              `existing`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// init
func init() {
	conf, viperConf = config.New()
	ctx = context.Background()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)

	// extra options that aren't handled via config env values
	// awscosts - month to get data for
	// existing - see if we want to add in costs
	existingCmd.Flags().BoolVar(&flagIncludeCosts, "include-costs", false, "When true, will impost cost data as well.")
}

func main() {
	rootCmd.AddCommand(existingCmd)
	err := rootCmd.Execute()
	// fail on errir
	if err != nil {
		os.Exit(1)
	}

}
