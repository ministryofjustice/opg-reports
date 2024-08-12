package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/shared/templ"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const defaultConfig string = "./config.json"

func main() {
	logger.LogSetup()
	ctx := context.Background()
	var err error
	var templates []string
	var configContent []byte

	mux := http.NewServeMux()
	// static assets
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	// favicon ignore
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// -- get templates
	templates = templ.GetTemplates("./templates")
	for _, f := range templates {
		slog.Debug("template file", slog.String("path", f))
	}

	// // -- get config
	configFile := env.Get("CONFIG_FILE", defaultConfig)
	if configContent, err = os.ReadFile(configFile); err != nil {
		slog.Error("error starting front...", slog.String("err", err.Error()))
		return
	}
	conf := config.New(configContent)

	// -- call github
	github_standards.Register(ctx, mux, conf, templates)

	addr := env.Get("FRONT_ADDR", ":8080")
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	slog.Info("starting web server",
		slog.String("log_level", logger.Level().String()),
		slog.String("address", addr),
	)
	server.ListenAndServe()
}