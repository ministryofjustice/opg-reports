package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/ministryofjustice/opg-reports/servers/front/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/front/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/template"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/govassets"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const defaultConfig string = "./config.json"
const defaultAddr string = ":8080"
const templateDir string = "./templates"

// download gov uk resources as a zip and extract
func init() {
	govassets.DownloadAssets()
}

func main() {
	logger.LogSetup()
	var (
		err           error
		templates     []string
		configContent []byte
		ctx           = context.Background()
		mux           = http.NewServeMux()
	)
	// handle static assets as directly from file system
	mux.Handle("/govuk/", http.StripPrefix("/govuk/", http.FileServer(http.Dir("govuk"))))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("govuk/assets"))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// favicon ignore
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// -- get templates
	templates = template.GetTemplates(templateDir)
	slog.Debug(convert.PrettyString(templates))

	// -- get config file loaded
	configFile := env.Get("CONFIG_FILE", defaultConfig)
	if configContent, err = os.ReadFile(configFile); err != nil {
		slog.Error("error reading config file", slog.String("err", err.Error()))
		return
	}
	conf := config.New(configContent)

	frontServer := front.New(ctx, conf, templates)
	// -- github standards
	github_standards.Register(mux, frontServer)
	// -- aws costs
	aws_costs.Register(mux, frontServer)

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
