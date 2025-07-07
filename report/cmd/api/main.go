package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/awsaccounts"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/awscosts"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/home"
	"github.com/ministryofjustice/opg-reports/report/cmd/api/teams"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/api"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
)

// root command
var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "API runner",
	Long:  `API enables the api server`,
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
