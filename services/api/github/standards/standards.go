package standards

import (
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/github/std"
	"opg-reports/shared/server"
	"opg-reports/shared/server/endpoint"
	"opg-reports/shared/server/resp"
)

var allowedParameters = []string{
	"archived",
	"team",
}

func Register(mux *http.ServeMux, store data.IStore[*std.Repository]) {

	qp := endpoint.NewQueryable(allowedParameters)

	mux.HandleFunc("/github/standards/{version}/list/{$}", func(w http.ResponseWriter, r *http.Request) {
		resp := resp.New()
		parameters := qp.Parse(r)
		filterFuncs := FilterFunctions(parameters, resp)

		data := endpoint.NewEndpointData[*std.Repository](store, nil, filterFuncs)
		display := endpoint.NewEndpointDisplay[*std.Repository](nil, DisplayRow, nil)
		ep := endpoint.New[*std.Repository]("test", resp, data, display, parameters)

		server.Middleware(ep.ProcessRequest, server.LoggingMW, server.SecurityHeadersMW)(w, r)
	})

}
