package daily

import (
	"log/slog"
	"net/http"
	"opg-reports/shared/aws/uptime"
	"opg-reports/shared/data"
	"opg-reports/shared/server/endpoint"
	"opg-reports/shared/server/mw"
	"opg-reports/shared/server/resp"
)

const taxServiceName string = "tax"

var allowedParameters = []string{
	"start",
	"end",
	"version",
}

func Register(mux *http.ServeMux, store data.IStore[*uptime.Uptime]) {
	slog.Info("registering routes",
		slog.String("handler", "aws-uptime-daily"),
		slog.Int("datastore count", store.Length()))
	qp := endpoint.NewQueryable(allowedParameters)

	// Note: if {start} or {end} are "-" it uses current month
	//
	// 	- /aws/uptime/{version}/daily/{start}/{end}/{$}
	mux.HandleFunc("/aws/uptime/{version}/daily/{start}/{end}/{$}", func(w http.ResponseWriter, r *http.Request) {
		response := resp.New()
		key := "monthly"

		parameters := qp.Parse(r)
		filterFuncs := FilterFunctions(parameters, response)
		displayHeadFuncs := DisplayHeadFunctions(parameters)
		displayRowFuncs := DisplayRowFunctions(parameters)

		head := displayHeadFuncs[key]
		row := displayRowFuncs[key]

		data := endpoint.NewEndpointData[*uptime.Uptime](store, nil, filterFuncs)
		display := endpoint.NewEndpointDisplay[*uptime.Uptime](head, row, nil)
		ep := endpoint.New[*uptime.Uptime]("aws-uptime-daily", response, data, display, parameters)

		mw.Middleware(ep.ProcessRequest, mw.Logging, mw.SecurityHeaders)(w, r)
	})

	mux.HandleFunc("/aws/uptime/{version}/daily/{start}/{end}/unit/{$}", func(w http.ResponseWriter, r *http.Request) {
		response := resp.New()
		key := "monthlyByAccountUnit"

		parameters := qp.Parse(r)
		filterFuncs := FilterFunctions(parameters, response)
		displayHeadFuncs := DisplayHeadFunctions(parameters)
		displayRowFuncs := DisplayRowFunctions(parameters)

		head := displayHeadFuncs[key]
		row := displayRowFuncs[key]

		data := endpoint.NewEndpointData[*uptime.Uptime](store, byUnit, filterFuncs)
		display := endpoint.NewEndpointDisplay[*uptime.Uptime](head, row, nil)
		ep := endpoint.New[*uptime.Uptime]("aws-uptime-daily-unit", response, data, display, parameters)

		mw.Middleware(ep.ProcessRequest, mw.Logging, mw.SecurityHeaders)(w, r)
	})

}
