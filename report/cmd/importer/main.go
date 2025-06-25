/*
Import data into a database

	importer
*/
package main

import (
	"github.com/ministryofjustice/opg-reports/report/cmd/importer/teams"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/spf13/cobra"
)

// root command
var rootCmd = &cobra.Command{
	Use:   "importer",
	Short: "Importer",
	Long:  `importer sub-commands will import data directly to the local database.`,
}

var conf, viperConf = config.New() // Get the configuration data and the viper config for mapping to cli args

// list of commands to attach to the root
var (
	teamCmd = teams.Cmd(conf, viperConf)
)

// init
// - Binds global config values to parameters
func init() {
	// bind the database.path config item to the shorter --db
	rootCmd.PersistentFlags().StringVar(&conf.Database.Path, "database.path", conf.Database.Path, "Path to database file")
	viperConf.BindPFlag("database.path", rootCmd.PersistentFlags().Lookup("database.path"))

}

func main() {

	rootCmd.AddCommand(teamCmd)
	rootCmd.Execute()

}
