package main

import (
	"context"
	"log/slog"

	"opg-reports/report/config"
	"opg-reports/report/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	conf      *config.Config
	viperConf *viper.Viper
	ctx       context.Context
	log       *slog.Logger
)

// root command
var rootCmd = &cobra.Command{
	Use:   "front",
	Short: "front end runner",
	Long: `
front turns on the front end server.

env values that can be adjusted:

	SERVERS_API_ADDR
		The address of the API server (eg: localhost:8081)
.
`,
	Run: func(cmd *cobra.Command, args []string) {
		runner(ctx, log, conf)
	},
}

// init called to setup conf / logger and to check and download
// govuk front end
func init() {
	conf, viperConf = config.New()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)
	ctx = context.Background()

	if !utils.DirExists(conf.GovUK.Front.Directory) {
		DownloadGovUKFrontEnd(ctx, log, conf)
	}
}

func main() {
	rootCmd.Execute()
}
