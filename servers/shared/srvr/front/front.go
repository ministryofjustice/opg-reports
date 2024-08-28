package front

import (
	"context"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/env"
)

type FrontServer struct {
	Ctx       context.Context
	Config    *config.Config
	Templates []string
	ApiSchema string
	ApiAddr   string
}

type FrontHandlerFunc func(server *FrontServer, navItem *nav.Nav, w http.ResponseWriter, r *http.Request)

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
