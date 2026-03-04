package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/account/accountapi/accountapi"
	"opg-reports/report/internal/codebases/codebasesapi/codebasesapistats"
	"opg-reports/report/internal/codeowners/codeownersapi"
	"opg-reports/report/internal/cost/costapi/costapiaccount"
	"opg-reports/report/internal/cost/costapi/costapidetailed"
	"opg-reports/report/internal/cost/costapi/costapidiff"
	"opg-reports/report/internal/cost/costapi/costapiteam"
	"opg-reports/report/internal/global/apimodels"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/internal/headline/headlineapi/headlineapi"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/internal/uptime/uptimeapi/uptimeapiteam"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/env"
	"opg-reports/report/package/logger"
	"os"

	"github.com/spf13/cobra"
)

type cli struct {
	DB      string `json:"db"`      // --db
	Driver  string `json:"driver"`  // --driver
	Params  string `json:"params"`  // --params
	ApiHost string `json:"api"`     // --api-host ; this is the server address to run from
	Version string `json:"version"` // --version ; the semver tag, used as part of signature
	SHA     string `json:"sha"`     // --sha ; the git commit sha used as part of signature
}

// default values for the args
var flags = &cli{
	Driver:  "sqlite3",
	DB:      "./database/api.db",
	ApiHost: ":8081",
	Version: "v0.0.0",
	SHA:     "abcde",
}

// main root command
var root *cobra.Command = &cobra.Command{
	Use:   "api",
	Short: `start the api`,
	RunE:  runAPI,
}

// registerEndpoints attaches all the current api endpoints into the
// server mux by calling the packages .Register function
//
// Adds a handler for / & /ping
func registerEndpoints(ctx context.Context, mux *http.ServeMux, in *cli) {

	var args = &apimodels.Args{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	}

	registerPingAndHome(ctx, mux, in)
	// teams
	// - all
	teamapiall.Register(ctx, mux, args)
	// headline api which returns highlighted numbers / optional team filter
	headlineapi.Register(ctx, mux, args)
	// accounts
	// - all accounts and team filter
	accountapi.Register(ctx, mux, args)
	// costs
	// - grouped by month & team / optional team filter
	costapiteam.Register(ctx, mux, args)
	// - grouped by account id / optional team filters
	costapiaccount.Register(ctx, mux, args)
	// - detailed costs group by month / optional team filter
	costapidetailed.Register(ctx, mux, args)
	// - cost differences / optional team filter
	costapidiff.Register(ctx, mux, args)
	// uptime
	// - uptime grouped by team name / optional team filter
	uptimeapiteam.Register(ctx, mux, args)
	// codebases
	// - stats / optional team filter
	codebasesapistats.Register(ctx, mux, args)
	// - ownership / optional team filter
	codeownersapi.Register(ctx, mux, args)
}

// runAPI the main run command
func runAPI(cmd *cobra.Command, args []string) (err error) {
	var (
		mux    *http.ServeMux
		server *http.Server
		ctx    context.Context = cmd.Context()
		log    *slog.Logger    = cntxt.GetLogger(ctx)
	)
	// overwrite arg flags from env values
	if err = env.OverwriteStruct(&flags); err != nil {
		return
	}
	// run db migrations
	err = migrations.Migrate(ctx, &migrations.Args{
		DB:     flags.DB,
		Driver: flags.Driver,
		Params: flags.Params,
	})
	if err != nil {
		return
	}

	// setup mux & server
	mux = http.NewServeMux()
	server = &http.Server{Addr: flags.ApiHost, Handler: mux}
	// attach endpoints
	registerEndpoints(ctx, mux, flags)
	// server info
	log.Info(fmt.Sprintf("Starting server [%s] [%s]...", flags.Version, flags.SHA))
	log.Info("Database:")
	log.Info(fmt.Sprintf("Driver: %s", flags.Driver))
	log.Info(fmt.Sprintf("Path: %s", flags.DB))
	log.Info("Hosts:")
	log.Info(fmt.Sprintf("API: [http://%s/]", flags.ApiHost))
	// boot server
	server.ListenAndServe()
	return
}

func init() {
	root.PersistentFlags().StringVar(&flags.Driver, "driver", flags.Driver, "Database driver")
	root.PersistentFlags().StringVar(&flags.DB, "db", flags.DB, "Database path")
	root.PersistentFlags().StringVar(&flags.Params, "params", flags.Params, "Database params")
	root.PersistentFlags().StringVar(&flags.ApiHost, "api-host", flags.ApiHost, "Address to run this api from")
	root.PersistentFlags().StringVar(&flags.Version, "version", flags.Version, "The semver")
	root.PersistentFlags().StringVar(&flags.SHA, "sha", flags.SHA, "The git commit sha")
}

func main() {
	var err error
	var log = logger.New(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_TYPE"))
	var ctx = cntxt.AddLogger(context.Background(), log)

	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error("error running command", "err", err.Error())
		os.Exit(1)
	}
}
