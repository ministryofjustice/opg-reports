package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/dl"
	"github.com/ministryofjustice/opg-reports/servers/front/front_templates"
	"github.com/ministryofjustice/opg-reports/servers/front/github_standards"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const defaultConfig string = "./config.json"
const defaultAddr string = ":8080"

// download gov uk resources as a zip and extract
func init() {
	dl.DownloadGovUKAssets()
}

func main() {
	fmt.Println("front running")
	logger.LogSetup()
	var ctx = context.Background()
	var err error
	var templates []string
	var configContent []byte
	var mux = http.NewServeMux()

	// handle static assets as directly from file system
	mux.Handle("/govuk/", http.StripPrefix("/govuk/", http.FileServer(http.Dir("govuk"))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// favicon ignore
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// -- get templates

	templates = front_templates.GetTemplates("./templates")
	for _, f := range templates {
		slog.Debug("template file", slog.String("path", f))
	}

	// -- get config
	configFile := env.Get("CONFIG_FILE", defaultConfig)
	if configContent, err = os.ReadFile(configFile); err != nil {
		slog.Error("error starting front - config file", slog.String("err", err.Error()))
		return
	}
	conf := config.New(configContent)

	// -- call github
	github_standards.Register(ctx, mux, conf, templates)

	addr := env.Get("FRONT_ADDR", defaultAddr)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	slog.Info("starting front server",
		slog.String("log_level", logger.Level().String()),
		slog.String("address", addr),
	)
	server.ListenAndServe()
}
