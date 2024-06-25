package main

import (
	"embed"
	"log/slog"
	"net/http"
	"opg-reports/services/front/cnf"
	"opg-reports/services/front/server"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/env"
	"opg-reports/shared/files"
	"opg-reports/shared/logger"
)

//go:embed templates
var templateFs embed.FS

//go:embed config.json
var configContent []byte

func main() {

	logger.LogSetup()

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	templateFS := files.NewFS(templateFs, "templates")
	templateFiles := tmpl.Files(templateFS, "")

	conf, _ := cnf.Load(configContent)
	serve := server.New(conf, templateFiles)
	serve.Register(mux)

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
