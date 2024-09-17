package front

import (
	"context"
	"log/slog"
	"net/http"
	"slices"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/mw"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/env"
)

type FrontHandlerFunc func(server *FrontServer, navItem *nav.Nav, w http.ResponseWriter, r *http.Request)

type FrontServer struct {
	Ctx        context.Context
	Config     *config.Config
	Templates  []string
	ApiSchema  string
	ApiAddr    string
	Registered []string
}

func (fs *FrontServer) HasHomePage() (hasHome bool) {
	hasHome = slices.Contains(fs.Registered, "/{$}")
	return
}

// Register func deals with adding url handling on the mux and tracks urls
func (fs *FrontServer) Register(mux *http.ServeMux, navItem *nav.Nav, handler FrontHandlerFunc) {
	var (
		httpHandler http.HandlerFunc = Wrap(fs, navItem, handler)
		uri         string           = navItem.Uri + "{$}"
	)
	if !slices.Contains(fs.Registered, uri) {
		// add to list
		fs.Registered = append(fs.Registered, uri)
		// add to mux
		mux.HandleFunc(uri, mw.Middleware(httpHandler, mw.Logging, mw.SecurityHeaders))
		// output info
		slog.Info("[front] registered", slog.String("uri", uri))
	}

}

// Wrap wraps a known function in a HandlerFunc and passes along the server and item details it needs
// outside of the normal scope of the http request
// Ths resulting function is then passed to the mux handle func (or via middleware)
func Wrap(server *FrontServer, navItem *nav.Nav, innerFunc FrontHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		innerFunc(server, navItem, w, r)
	}
}

func New(ctx context.Context, conf *config.Config, templates []string) *FrontServer {
	return &FrontServer{
		Ctx:       ctx,
		Config:    conf,
		Templates: templates,
		ApiSchema: env.Get("API_SCHEME", consts.API_SCHEME),
		ApiAddr:   env.Get("API_ADDR", consts.API_ADDR),
	}
}
