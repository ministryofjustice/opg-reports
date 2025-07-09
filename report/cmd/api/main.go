package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	conf      *config.Config
	viperConf *viper.Viper
	ctx       context.Context
	log       *slog.Logger
)

var maxDatabaseAge time.Duration = (24 * time.Hour) * 3 // max age of the database before refetching - 3 days

// root command
var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "API runner",
	Long: `
api turns on the api and the docs - generated using huma.

env values that can be adjusted:

	AWS_BUCKETS_DB_NAME
		The name of the bucket where latest database is stored
	AWS_BUCKETS_DB_KEY
		The object key in the bucket (including folder path) for the latest database
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

Requires valid AWS session with permission to access AWS_BUCKETS_DB_NAME and AWS_BUCKETS_DB_KEY.
`,
	Run: func(cmd *cobra.Command, args []string) {
		runner(ctx, log, conf)
	},
}

// init called to setup conf / logger and to check database exists or create it
func init() {
	conf, viperConf = config.New()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)
	ctx = context.Background()

	if !utils.FileExists(conf.Database.Path) {
		seedDB(ctx, log, conf)
	}
}

func main() {
	rootCmd.Execute()
}
