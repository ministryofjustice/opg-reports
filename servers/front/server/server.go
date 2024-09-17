package server

import (
	"context"
	"log/slog"
	"net/http"
	"slices"

	"github.com/ministryofjustice/opg-reports/servers/front/apihandlers"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/template"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/mw"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/env"
)

type FrontServer struct {
	Ctx        context.Context
	Config     *config.Config
	ApiSchema  string
	ApiAddr    string
	Mux        *http.ServeMux
	Templates  []string
	Registered []string
	PageData   map[string]interface{}
}

func (svr *FrontServer) RegisterPage(uri string, navItem *nav.Nav) {
	suffix := "{$}"
	svr.Registered = append(svr.Registered, uri)
	frontHandler := wrap(svr, navItem)
	slog.Info("[front] registering", slog.String("uri", uri))
	svr.Mux.HandleFunc(uri+suffix, mw.Middleware(frontHandler, mw.Logging, mw.SecurityHeaders))
}

func (svr *FrontServer) Register() {
	statics(svr)
	allNav := nav.Flattern(svr.Config.Navigation)
	for _, navItem := range allNav {
		if navItem.Uri != "" {
			svr.RegisterPage(navItem.Uri, navItem)
		}
	}
}

func (svr *FrontServer) HasHomepage() bool {
	return slices.Contains(svr.Registered, "/")
}

func New(ctx context.Context, mux *http.ServeMux, cfg *config.Config, templates []string) (svr *FrontServer) {
	svr = &FrontServer{
		Ctx:        ctx,
		Config:     cfg,
		Templates:  templates,
		Mux:        mux,
		Registered: []string{},
		PageData:   map[string]interface{}{},
		ApiSchema:  env.Get("API_SCHEME", consts.API_SCHEME),
		ApiAddr:    env.Get("API_ADDR", consts.API_ADDR),
	}

	return
}

func statics(svr *FrontServer) {
	// handle static assets as directly from file system
	svr.Mux.Handle("/govuk/", http.StripPrefix("/govuk/", http.FileServer(http.Dir("govuk"))))
	svr.Mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("govuk/assets"))))
	svr.Mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// favicon ignore
	svr.Mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

}

func generatePage(svr *FrontServer, navItem *nav.Nav, writer http.ResponseWriter, request *http.Request) {
	pageTemplate := template.New(navItem.Template, svr.Templates, writer)

	for key, path := range navItem.DataSources {
		if handler := apihandlers.Get(path); handler != nil {
			remote := httphandler.New(svr.ApiSchema, svr.ApiAddr, path)
			handler.Handle(svr.PageData, key, remote, writer, request)
		}
	}
	pageData := NewPage(svr, navItem, request)
	pageTemplate.Run(pageData)
}

// wrap wraps a known function in a HandlerFunc and passes along the server and item details it needs
// outside of the normal scope of the http request
// Ths resulting function is then passed to the mux handle func (or via middleware)
func wrap(server *FrontServer, navItem *nav.Nav) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		generatePage(server, navItem, w, r)
	}
}
