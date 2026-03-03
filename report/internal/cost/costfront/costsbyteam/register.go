package costsbyteam

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/report/internal/global/frontmodels"
	"opg-reports/report/package/cntxt"
)

const ENDPOINT string = `/home/costs/teams/`

var endpoints []string = []string{
	ENDPOINT,
}

func Register(ctx context.Context, mux *http.ServeMux, args *frontmodels.RegisterArgs) {
	var log = cntxt.GetLogger(ctx)

	for _, ep := range endpoints {
		log.Info(fmt.Sprintf("[%s] registering endpoint [%s] to handler", "costsbyteam", ep))
		ep = fmt.Sprintf("%s{$}", ep)

		mux.HandleFunc(ep, func(writer http.ResponseWriter, request *http.Request) {
			Handler(ctx, args, request, writer)
		})
	}
}
