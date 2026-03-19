package home

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/report/packages/httpx"
	"opg-reports/report/packages/slogx"
)

// Assets maps the `/assets/` path
//
// `/assets/` is hardcorded within the govuk css and js for where fonts / images are, so map that
// to the filesystem (./govuk/assets/)
//
// Example:
//   - http://localhost:8080/assets/images/govuk-icon-180.png => ./govuk/assets/images/govuk-icon-180.png
func Assets(ctx context.Context, mux httpx.Mux, muxcfg httpx.MuxConfig) {
	var log = slogx.FromContext(ctx)
	var cfg = muxcfg.Conf()
	var dir = cfg.GovUKAssetsDir()
	var ep = `/assets/`

	log.Info(ctx, fmt.Sprintf(`registering endpoint [%s] via mux handle from directory [%s]`, ep, dir))
	mux.Handle(ep, http.FileServer(http.Dir(dir)))
}

// GovUK handles the `/govuk/` path
//
// `/govuk/` is path we use to include css / js from gov uk assets, so capture and point to the gov uk directory
//
// Example:
//   - http://localhost:8080/govuk/govuk-frontend-5.14.0.min.css => ./govuk/govuk-frontend-5.14.0.min.css
func GovUK(ctx context.Context, mux httpx.Mux, muxcfg httpx.MuxConfig) {
	var log = slogx.FromContext(ctx)
	var cfg = muxcfg.Conf()
	var dir = cfg.GovUKAssetsDir()
	var ep = `/govuk/`

	log.Info(ctx, fmt.Sprintf(`registering endpoint [%s] via mux handle from directory [%s]`, ep, dir))
	mux.Handle(ep, http.StripPrefix(ep, http.FileServer(http.Dir(dir))))
}

// LocalAssets handles the `/local-assets/` path
//
// `/local-assets/` contain our css overwrites, extra images / js and so on
//
// Example:
//   - http://localhost:8080/local-assets/css/local.css => ./local-assets/css/local.css
func LocalAssets(ctx context.Context, mux httpx.Mux, muxcfg httpx.MuxConfig) {
	var log = slogx.FromContext(ctx)
	var cfg = muxcfg.Conf()
	var dir = cfg.LocalAssetsDir()
	var ep = `/local-assets/`

	log.Info(ctx, fmt.Sprintf(`registering endpoint [%s] via mux handle from directory [%s]`, ep, dir))
	mux.Handle(ep, http.StripPrefix(ep, http.FileServer(http.Dir(dir))))
}

// IgnoreFavicon
func IgnoreFavicon(ctx context.Context, mux httpx.Mux, muxcfg httpx.MuxConfig) {
	var log = slogx.FromContext(ctx)
	var ep = `/favicon.ico`

	log.Info(ctx, fmt.Sprintf(`registering endpoint [%s] to ignore favicons`, ep))
	mux.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Not found.`))
	})
}
