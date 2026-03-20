package main

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/report/internal/config"
	"opg-reports/report/internal/migrations"
	"opg-reports/report/packages/httpx"
	"opg-reports/report/packages/slogx"
	"os"

	"github.com/spf13/cobra"
)

var cfg *config.Config // the main api config

// main root command
var root *cobra.Command = &cobra.Command{
	Use:   "api",
	Short: `api starts the main api handler command.`,
	RunE:  runCmd,
}

// runCmd is the main cobra command
func runCmd(cmd *cobra.Command, args []string) (err error) {
	var (
		server *http.Server
		ctx    context.Context = cmd.Context()
		log    slogx.Logger    = slogx.FromContext(ctx)
		mux    httpx.MuxServer = httpx.NewMux()
	)
	// run migrations
	err = migrations.Migrate(ctx, cfg)
	if err != nil {
		log.Error(ctx, "error running migrations", "err", err.Error())
	}

	// create the server
	server = &http.Server{
		Addr:    cfg.ApiHostname(),
		Handler: mux,
	}

	// attach endpoints
	registerEndpoints(ctx, mux, cfg)

	log.Info(ctx, `starting api server ...`,
		"database", cfg.DBPath,
		"api host", fmt.Sprintf(`http://%s`, cfg.ApiHostname()),
	)
	defer log.Info(ctx, `shutting down api server ...`)

	// start the server
	server.ListenAndServe()
	return
}

func init() {
	// setup the config and bind to the command
	cfg = config.NewApi()
	cfg.BindApi(root)
}

func main() {
	var (
		err error
		ctx context.Context = context.Background()
		log slogx.Logger    = slogx.New(slogx.Config())
	)
	// attach the log & config to the current context
	ctx = slogx.Attach(ctx, log)
	// run root command
	err = root.ExecuteContext(ctx)
	if err != nil {
		log.Error(ctx, "error running command", "err", err.Error())
		os.Exit(1)
	}
}
