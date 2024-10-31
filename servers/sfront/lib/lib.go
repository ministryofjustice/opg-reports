package lib

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ministryofjustice/opg-reports/pkg/govukassets"
	"github.com/ministryofjustice/opg-reports/pkg/navigation"
)

// DownloadGovUKAssets fetches the assets from govuk front end
// and moves it to directory
func DownloadGovUKAssets(directory string) (err error) {
	var frontEnd = govukassets.FrontEnd()
	defer frontEnd.Close()

	_, err = frontEnd.Do(directory)
	return

}

// TemplateFiles recurisvely finds all go templates (`.gotmpl`) within directory
func TemplateFiles(directory string) (templates []string) {
	var ext = ".gotmpl"
	templates = []string{}
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(info.Name()) == ext {
			templates = append(templates, path)
		}
		return nil
	})
	return
}

// Statics handles the top level static folder paths to use files and not active code as
// well as erroring for the default favicon
//
// - /assets/*
// - /static/*
func Statics(mux *http.ServeMux) {
	slog.Info("registering statics ...")

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}

// HomepageRedirect looks for a homepage, if it cant find one it will
// set a redirect for / to the first nav item
func HomepageRedirect(mux *http.ServeMux, flatNav map[string]*navigation.Navigation, first *navigation.Navigation) {
	var (
		home        = "/"
		suffix      = "{$}"
		destination = first.Uri
		redirector  = func(w http.ResponseWriter, r *http.Request) {
			slog.Info("[sfront] redirecting ...",
				slog.String("origin", r.URL.String()),
				slog.String("destination", destination))
			http.Redirect(w, r, destination, http.StatusTemporaryRedirect)
		}
	)

	if _, ok := flatNav[home]; !ok {
		mux.HandleFunc(home+suffix, redirector)
	}

}
