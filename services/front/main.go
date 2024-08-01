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
	"os"
)

//go:embed templates/**
var templateEmbed embed.FS

func main() {

	logger.LogSetup()

	mux := http.NewServeMux()
	// static assets
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	// favicon ignore
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	templateDir := os.DirFS("templates/").(files.IReadFS)
	templateFS := files.NewFS(templateDir, "templates/")
	templateFiles := tmpl.Files(templateFS, "templates/")

	for _, f := range templateFiles {
		slog.Debug("template file", slog.String("path", f))
	}
	configFile := env.Get("CONFIG_FILE", "./config.opg.json")
	configContent, _ := os.ReadFile(configFile)

	conf, _ := cnf.Load(configContent)
	serve := server.New(
		conf,
		templateFiles,
		env.Get("API_ADDR", ":8081"),
		env.Get("API_SCHEME", "http"),
	)

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
