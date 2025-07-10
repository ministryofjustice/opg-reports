package main

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
)

func RegisterStaticHandlers(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	mux *http.ServeMux,
) {
	log.Info("registering static handlers ...")
	// Static assets
	// /assets/ is hardcorded in the govuk css and js for where fonts / images are, so map that to the filesystem (./govuk/assets/)
	// 		http://localhost:8080/assets/images/govuk-icon-180.png
	mux.Handle("/assets/", http.FileServer(http.Dir(govUKAssetDir)))
	// /local-assets/ contain our css overwrites, extra images / js and so on
	//		http://localhost:8080/local-assets/css/local.css
	mux.Handle("/local-assets/", http.StripPrefix("/local-assets/", http.FileServer(http.Dir(localAssetsDir))))
	// /govuk/ is path we use to include css / js, so capture and point to the gov uk directory
	// 		http://localhost:8080/govuk/VERSION.TXT
	// 		http://localhost:8080/govuk/govuk-frontend-5.11.0.min.css
	mux.Handle("/govuk/", http.StripPrefix("/govuk/", http.FileServer(http.Dir(govUKAssetDir))))
	// ignore favicons
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

}
