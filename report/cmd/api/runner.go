package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"opg-reports/report/config"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
)

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
