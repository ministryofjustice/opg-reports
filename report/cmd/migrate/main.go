/*
migrate used to run the schema to migrate db

# Example

`env DATABASE_PATH="./data/api.db" migrate up`
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

// root command
var rootCmd = &cobra.Command{
	Use:               "migrate",
	Short:             "migrate",
	Long:              `migrate the database`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// init
func init() {
	conf, viperConf = config.New()
	ctx = context.Background()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)
}

func main() {
	rootCmd.AddCommand(upCmd)
	err := rootCmd.Execute()
	// fail on errir
	if err != nil {
		os.Exit(1)
	}

}
