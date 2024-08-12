package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/shared/templ"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const defaultConfig string = "./config.json"

const govukVersion string = "v5.4.0"

// https://github.com/alphagov/govuk-frontend/releases/download/v5.4.0/release-v5.4.0.zip

// download gov uk resources as a zip and extract
func init() {
	slog.Warn("Downloading govuk assets", slog.String("v", govukVersion))
	// url for the public download
	zipUrl := fmt.Sprintf("https://github.com/alphagov/govuk-frontend/releases/download/%s/release-%s.zip",
		govukVersion, govukVersion)

	resp, err := http.Get(zipUrl)
	if err != nil {
		slog.Error("error getting assets", slog.String("err", err.Error()))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("error status with assets", slog.Int("status", resp.StatusCode))
		return
	}
	os.MkdirAll("govuk", os.ModePerm)
	// Create the file
	zFile := "govuk/assets.zip"
	out, err := os.Create(zFile)
	if err != nil {
		slog.Error("error creating assets zip", slog.String("err", err.Error()))
		return
	}

	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		slog.Error("error writing to zip", slog.String("err", err.Error()))
		return
	}
	// extract
	archive, err := zip.OpenReader(zFile)
	defer archive.Close()
	if err != nil {
		slog.Error("error opening zip", slog.String("err", err.Error()))
		return
	}

	for _, f := range archive.File {
		path := filepath.Join("govuk", f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			// make parent dir
			os.MkdirAll(filepath.Dir(path), os.ModePerm)
			dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			defer dstFile.Close()
			if err != nil {
				slog.Error("error opening file", slog.String("err", err.Error()))
				return
			}
			srcFile, err := f.Open()
			defer srcFile.Close()
			if err != nil {
				slog.Error("error opening file", slog.String("err", err.Error()))
				return
			}
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				slog.Error("error opening file", slog.String("err", err.Error()))
				return
			}

		}
	}
	// remove the zip
	os.Remove(zFile)

}

func main() {
	fmt.Println("main")
	logger.LogSetup()
	ctx := context.Background()
	var err error
	var templates []string
	var configContent []byte

	mux := http.NewServeMux()
	// static assets
	mux.Handle("/govuk/", http.StripPrefix("/govuk/", http.FileServer(http.Dir("govuk"))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
