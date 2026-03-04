package codeownersfront

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/global/frontmodels"
	"opg-reports/report/package/cntxt"
)

const ENDPOINT_BASE string = `/home/codeownership/`
const ENDPOINT_TEAM string = `/team/{team}/codeownership/`

var endpoints []string = []string{
	ENDPOINT_BASE,
	ENDPOINT_TEAM,
}

func Register(ctx context.Context, mux *http.ServeMux, args *frontmodels.RegisterArgs) {
	var log *slog.Logger = cntxt.GetLogger(ctx)

	for _, ep := range endpoints {
		log.Info(fmt.Sprintf("[%s] registering endpoint [%s] to handler", "codeownersfront", ep))
		ep = fmt.Sprintf("%s{$}", ep)

		mux.HandleFunc(ep, func(writer http.ResponseWriter, request *http.Request) {
			Handler(ctx, args, request, writer)
		})
	}
}
