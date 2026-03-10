package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/migrations"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/logger"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var flags *args.API

// main root command
var root *cobra.Command = &cobra.Command{
	Use:   "api",
	Short: `run the api`,
	RunE:  runAPI,
}

// runAPI the main run command
func runAPI(cmd *cobra.Command, args []string) (err error) {
	var (
		log    *slog.Logger
		mux    *http.ServeMux
		server *http.Server
		ctx    context.Context = cmd.Context()
	)
	ctx, log = logger.Get(ctx)

	// // overwrite arg flags from env values
	// if err = env.OverwriteStruct(&flags); err != nil {
	// 	return
	// }

	// run db migrations
	err = migrations.Migrate(ctx, flags.DB)
	if err != nil {
		return
	}

	// setup mux & server
	mux = http.NewServeMux()
	server = &http.Server{Addr: flags.Hosts.API, Handler: mux}
	// attach endpoints
	registerEndpoints(ctx, mux, flags)
	// server info
	log.Info(fmt.Sprintf("Starting server [%s] [%s]...", flags.Versions.Version, flags.Versions.SHA))
	log.Info("Database:")
	log.Info(fmt.Sprintf("Driver: %s", flags.DB.Driver))
	log.Info(fmt.Sprintf("Path: %s", flags.DB.DB))
	log.Info("Hosts:")
	log.Info(fmt.Sprintf("API: [http://%s/]", flags.Hosts.API))
	// boot server
	server.ListenAndServe()
	return
}

func defaultArgs(t time.Time) (def *args.API) {
	def = args.Default[*args.API](t)
	return
}

func init() {
	var now = time.Now().UTC()

	flags = defaultArgs(now)
	// setup all the flags
	// - DB
	root.PersistentFlags().StringVar(&flags.DB.Driver, "driver", flags.DB.Driver, "database driver type.")
	root.PersistentFlags().StringVar(&flags.DB.DB, "db", flags.DB.DB, "database path.")
	root.PersistentFlags().StringVar(&flags.DB.Params, "params", flags.DB.Params, "database connection parameters.")
	// - Versions
	root.PersistentFlags().StringVar(&flags.Versions.Version, "version", flags.Versions.Version, "semver version.")
	root.PersistentFlags().StringVar(&flags.Versions.SHA, "sha", flags.Versions.SHA, "git sha.")
	// - hosts
	root.PersistentFlags().StringVar(&flags.Hosts.API, "api-host", flags.Hosts.API, "api host address.")

}

func main() {
	var err error
	var log *slog.Logger
	var ctx = context.Background()
	ctx, log = logger.Get(ctx)

	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error("error with command", "err", err.Error())
		panic("error")
		os.Exit(1)
	}

}
