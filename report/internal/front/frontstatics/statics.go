package frontstatics

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/package/cntxt"
	"strings"
)

type Args struct {
	RootDir        string `json:"root_dir"`
	GovUKDir       string `json:"govuk_dir"`
	LocalAssetsDir string `json:"local_assets_dir"`
}

// Register setups up how the server handles govuk
// assets and local css / image files and manages mapping of the
// urls to folders.
//
// Required as the govuk css / js contains references to resources
// like fonts & images with fixed paths (generally /assets/) which
// does not match our folder structure after zip extraction.
//
// Example url to filesytems mapping:
//
//	http://localhost:8080/assets/images/govuk-icon-180.png 		=> ./govuk/assets/images/govuk-icon-180.png
//	http://localhost:8080/local-assets/css/local.css 			=> ./local-assets/css/local.css
//	http://localhost:8080/govuk/govuk-frontend-5.11.0.min.css 	=> ./govuk/govuk-frontend-5.11.0.min.css
func Register(ctx context.Context, mux *http.ServeMux, dirs *Args) {
	var prefix = ""
	var root = strings.TrimPrefix(dirs.RootDir, "./")
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "frontstatics", "func", "Register")
	// Static assets
	// /assets/ is hardcorded in the govuk css and js for where fonts / images are, so map that to the filesystem (./govuk/assets/)
	// 		http://localhost:8080/assets/images/govuk-icon-180.png
	log.Info("registering handler [`/assets/`] ...")
	mux.Handle("/assets/", http.FileServer(http.Dir(dirs.GovUKDir)))

	// /local-assets/ contain our css overwrites, extra images / js and so on
	//		http://localhost:8080/local-assets/css/local.css
	prefix = strings.ReplaceAll(dirs.LocalAssetsDir, root, "") + "/"
	log.Info("registering handler [`" + prefix + "`] ...")
	mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(dirs.LocalAssetsDir))))

	// /govuk/ is path we use to include css / js, so capture and point to the gov uk directory
	// 		http://localhost:8080/govuk/VERSION.TXT
	// 		http://localhost:8080/govuk/govuk-frontend-5.11.0.min.css
	prefix = strings.ReplaceAll(dirs.GovUKDir, root, "") + "/"
	log.Info("registering handler [`" + prefix + "`] ...")
	mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(dirs.GovUKDir))))

	// ignore favicons
	log.Info("registering handler [`/favicon.ico`] ...")
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

}
