package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/awsaccounts"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/awscosts"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/home"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/teams"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/api"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
)

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
		var (
			ctx context.Context = context.Background()
			log *slog.Logger    = utils.Logger(conf.Log.Level, conf.Log.Type)
		)
		runner(ctx, log, conf)
	},
}

// Get the configuration data and the viper config for mapping to cli args
var conf, viperConf = config.New()

// addMiddleware add all middleware into the reqquest; currently these are:
//
//   - Check max age of the local database, if older than 3 days, fetch from s3
func addMiddleware(hapi huma.API, log *slog.Logger, conf *config.Config) {
	// add database age information
	hapi.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		var (
			err          error
			stats        fs.FileInfo
			modifiedTime time.Time
			age          time.Duration = 0 * time.Second
			now          time.Time     = time.Now().UTC()
			path         string        = conf.Database.Path
			maxAge       time.Duration = (24 * time.Hour) * 3 // 3 days
			// maxAge time.Duration = 10 * time.Minute
		)

		if stats, err = os.Stat(path); err == nil {
			modifiedTime = stats.ModTime()
			age = now.Sub(modifiedTime)
		}
		// if the age is above the max allowed age, fetch the database
		if age > maxAge {
			log.Warn("database age is beyond max ... downloading ...", "age", age, "destination", path)
			if err = downloadLatestDB(ctx.Context(), log, conf); err != nil {
				log.Error("error downloading updated database", "err", err.Error())
			} else {
				log.Info("downloaded new database from s3 ...")
			}
		} else {
			log.Info("database age ... ", "age", age)
		}

		next(ctx)
	})

}

// downloadLatestDB uses the config settings to create a s3 client and service to then download
// the latest database from the s3 bucket to local files and then updates the modified times.
//
// Changing the modified time is done to ensure fresh database is constantly being pulled from
// s3 when the database is older (say db has not been updated in 5 days).
func downloadLatestDB(ctx context.Context, log *slog.Logger, conf *config.Config) (err error) {
	var (
		local  string
		client = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
		store  = awsr.Default(ctx, log, conf)
		dir, _ = os.MkdirTemp("./", "__download-s3-*")
		now    = time.Now().UTC()
	)
	defer os.RemoveAll(dir)
	local, err = store.DownloadItemFromBucket(client, conf.Aws.Buckets.DB.Name, conf.Aws.Buckets.DB.Path(), dir)
	if err != nil {
		return
	}
	err = os.Rename(local, conf.Database.Path)
	if err != nil {
		return
	}
	err = os.Chtimes(conf.Database.Path, now, now)

	return
}

// RegisterHandlers attaches all the known functions to the api.
//
// To allow for service injection, each is called directly, so need to be manually added
func RegisterHandlers(ctx context.Context, log *slog.Logger, conf *config.Config, humaapi huma.API) {
	var (
		teamStore          = sqlr.DefaultWithSelect[*api.Team](ctx, log, conf)
		teamService        = api.Default[*api.Team](ctx, log, conf)
		awsAccountsStore   = sqlr.DefaultWithSelect[*api.AwsAccount](ctx, log, conf)
		awsAccountService  = api.Default[*api.AwsAccount](ctx, log, conf)
		awsCostsStore      = sqlr.DefaultWithSelect[*api.AwsCost](ctx, log, conf)
		awsCostsService    = api.Default[*api.AwsCost](ctx, log, conf)
		awsCostsStoreGroup = sqlr.DefaultWithSelect[*api.AwsCostGrouped](ctx, log, conf)
		awsCostsSrvGroup   = api.Default[*api.AwsCostGrouped](ctx, log, conf)
	)
	// HOME
	home.RegisterGetHomepage(log, conf, humaapi)
	// TEAMS
	teams.RegisterGetTeamsAll(log, conf, humaapi, teamService, teamStore)
	// AWS ACCOUNTS
	awsaccounts.RegisterGetAwsAccountsAll(log, conf, humaapi, awsAccountService, awsAccountsStore)
	// AWS COSTS
	awscosts.RegisterGetAwsCostsTop20(log, conf, humaapi, awsCostsService, awsCostsStore)
	awscosts.RegisterGetAwsGroupedCosts(log, conf, humaapi, awsCostsSrvGroup, awsCostsStoreGroup)
}

func runner(ctx context.Context, log *slog.Logger, conf *config.Config) {
	var (
		humaapi       huma.API
		cli           humacli.CLI
		server        http.Server
		mux           *http.ServeMux = http.NewServeMux()
		apiName       string         = conf.Servers.Api.Name
		apiVersion    string         = fmt.Sprintf("%s [%s]", conf.Versions.Semver, conf.Versions.Commit)
		shutdownDelay time.Duration  = 5 * time.Second
	)

	// create the server
	server = http.Server{
		Addr:    conf.Servers.Api.Addr,
		Handler: mux,
	}
	// create the api
	humaapi = humago.New(mux, huma.DefaultConfig(apiName, apiVersion))
	cli = humacli.New(func(hooks humacli.Hooks, opts *struct{}) {
		var addr = server.Addr

		// Inject middleware to api requests
		addMiddleware(humaapi, log, conf)

		RegisterHandlers(ctx, log, conf, humaapi)
		// startup
		hooks.OnStart(func() {
			log.Info("Starting api server...")
			log.Info(fmt.Sprintf("API: [http://%s/]", addr))
			log.Info(fmt.Sprintf("Docs: [http://%s/docs]", addr))

			server.ListenAndServe()
		})
		// graceful shutdown
		hooks.OnStop(func() {
			slog.Info("Stopping api server...")
			ctx, cancel := context.WithTimeout(ctx, shutdownDelay)
			defer cancel()
			server.Shutdown(ctx)
		})

	})
	cli.Run()
}

func main() {
	rootCmd.Execute()
}
