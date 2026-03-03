package costapidiff

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/report/internal/global/apimodels"
	"opg-reports/report/package/cntxt"
)

const ENDPOINT_BASE string = `/v1/costs/diff/{date_a}/{date_b}/`
const ENDPOINT_TEAM string = `/v1/costs/diff/{date_a}/{date_b}/team/{team}/`

var endpoints []string = []string{
	ENDPOINT_BASE,
	ENDPOINT_TEAM,
}

// Register wraps the handle func with a local version that also gets additional config
// details
func Register(ctx context.Context, mux *http.ServeMux, config *apimodels.Args) {
	var log = cntxt.GetLogger(ctx)

	for _, ep := range endpoints {
		log.Info(fmt.Sprintf("[%s] registering endpoint [%s] to handler", "costapidiff", ep))
		ep = fmt.Sprintf("%s{$}", ep)

		mux.HandleFunc(ep, func(writer http.ResponseWriter, request *http.Request) {
			Responder(ctx, config, request, writer)
		})
	}

}
