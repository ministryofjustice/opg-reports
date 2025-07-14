package main

import (
	"context"
	"log/slog"

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
	Use:               "db",
	Short:             "db",
	Long:              `db downloads or uploads sqlite database from s3`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// init
func init() {
	conf, viperConf = config.New()
	ctx = context.Background()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)

}

func main() {
	rootCmd.AddCommand(downloadCmd, uploadCmd)
	rootCmd.Execute()

}
