package teamapiall

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/report/internal/global/apimodels"
	"opg-reports/report/package/cntxt"
)

const ENDPOINT string = "/v1/teams"

// Register wraps the handle func with a local version that also gets additional config
// details
func Register(ctx context.Context, mux *http.ServeMux, config *apimodels.Args) {
	var log = cntxt.GetLogger(ctx)

	log.Info("registering handler ... ", "endpoint", ENDPOINT)
	mux.HandleFunc(fmt.Sprintf("%s/{$}", ENDPOINT), func(writer http.ResponseWriter, request *http.Request) {
		Responder(ctx, config, request, writer)
	})

}
