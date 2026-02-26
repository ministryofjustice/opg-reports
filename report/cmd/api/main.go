package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/account/accountapi/accountapiall"
	"opg-reports/report/internal/account/accountapi/accountsapiforteam"
	"opg-reports/report/internal/codebases/codebasesapi/codebasestatsapi"
	"opg-reports/report/internal/cost/costapi/costapidetailed"
	"opg-reports/report/internal/cost/costapi/costapidetailedteamfilter"
	"opg-reports/report/internal/cost/costapi/costapidiff"
	"opg-reports/report/internal/cost/costapi/costapidiffteamfilter"
	"opg-reports/report/internal/cost/costapi/costapiteam"
	"opg-reports/report/internal/cost/costapi/costapiteamfilter"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/internal/headline/headlineapi/headlineapihome"
	"opg-reports/report/internal/headline/headlineapi/headlineapiteam"
	"opg-reports/report/internal/team/teamapi/teamapiall"
	"opg-reports/report/internal/uptime/uptimeapi/uptimeapiteam"
	"opg-reports/report/internal/uptime/uptimeapi/uptimeapiteamfilter"
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

	registerPingAndHome(ctx, mux, in)
	// teams
	// - all
	teamapiall.Register(ctx, mux, &teamapiall.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// headline api which returns highlighted numbers
	headlineapihome.Register(ctx, mux, &headlineapihome.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// headline figures for team pages
	headlineapiteam.Register(ctx, mux, &headlineapiteam.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})

	// accounts
	// - all
	accountapiall.Register(ctx, mux, &accountapiall.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// - filtered by a team
	accountsapiforteam.Register(ctx, mux, &accountsapiforteam.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})

	// costs
	// - grouped by month & team
	costapiteam.Register(ctx, mux, &costapiteam.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// - grouped by month, filtered by team
	costapiteamfilter.Register(ctx, mux, &costapiteamfilter.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// - detailed costs group by month
	costapidetailed.Register(ctx, mux, &costapidetailed.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// - detailed view grouped by month - filtered by team
	costapidetailedteamfilter.Register(ctx, mux, &costapidetailedteamfilter.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// - cost differences
	costapidiff.Register(ctx, mux, &costapidiff.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// - cost difference filtered by a team
	costapidiffteamfilter.Register(ctx, mux, &costapidiffteamfilter.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})

	// uptime
	// - uptime grouped by team name
	uptimeapiteam.Register(ctx, mux, &uptimeapiteam.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// - uptime filtered by team
	uptimeapiteamfilter.Register(ctx, mux, &uptimeapiteamfilter.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
	// codebases
	codebasestatsapi.Register(ctx, mux, &codebasestatsapi.Config{
		DB:      in.DB,
		Driver:  in.Driver,
		Params:  in.Params,
		Version: in.Version,
		SHA:     in.SHA,
	})
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
