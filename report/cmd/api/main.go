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

var (
	conf      *config.Config
	viperConf *viper.Viper
	ctx       context.Context
	log       *slog.Logger
)

// root command
var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "API runner",
	Long: `
api turns on the api and the docs - generated using huma.

env values that can be adjusted:

	DATABASE_PATH
		The local file system path to the database the API is using
	SERVERS_API_NAME
		The name this API uses in docs
	SERVERS_API_ADDR
		The address of the API server (eg: localhost:8081)
	VERSIONS_SEMVER
		The semantic version tag this API is built from (eg: v3.0.1)
	VERSIONS_COMMIT
		The git commit hash that was used to build this version of the API

Requires valid AWS session with permission to access DATABASE_BUCKET_NAME and DATABASE_BUCKET_KEY.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		runner(ctx, log, conf)
		return
	},
}

// init called to setup conf / logger and to check database exists or create it
func init() {
	conf, viperConf = config.New()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)
	ctx = context.Background()

	if !utils.FileExists(conf.Database.Path) {
		log.Error("database not found")
		os.Exit(1)
	}
}

func main() {
	err := rootCmd.Execute()
	// fail on errir
	if err != nil {
		os.Exit(1)
	}
}
