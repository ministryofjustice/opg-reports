package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/ministryofjustice/opg-reports/servers/front/server"
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
		frontServer   *server.FrontServer
		ctx           context.Context = context.Background()
		mux                           = http.NewServeMux()
	)

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
	frontServer = server.New(ctx, mux, conf, templates)
	frontServer.Register()

	// -- if home page is not set, then give a default
	if !frontServer.HasHomepage() {
		home := conf.Navigation[0]
		frontServer.RedirectPage("/", home.Uri)
	}

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
