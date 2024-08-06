package monthly

import (
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/server"
	"opg-reports/shared/server/endpoint"
	"opg-reports/shared/server/resp"
)

const taxServiceName string = "tax"

var allowedParameters = []string{
	"start",
	"end",
	"version",
	"unit",
	"environment",
}

func Register(mux *http.ServeMux, store data.IStore[*cost.Cost]) {

	qp := endpoint.NewQueryable(allowedParameters)

	// Returns cost data split into with & without tax segments, then grouped by the month
	// Previously "Totals" sheet
	// Note: if {start} or {end} are "-" it uses current month
	//
	// 	- /aws/costs/{version}/monthly/{start}/{end}/{$}
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/{$}", func(w http.ResponseWriter, r *http.Request) {
		response := resp.New()
		key := "monthlyTotals"

		parameters := qp.Parse(r)
		filterFuncs := FilterFunctions(parameters, response)
		displayHeadFuncs := DisplayHeadFunctions(parameters)
		displayRowFuncs := DisplayRowFunctions(parameters)

		head := displayHeadFuncs[key]
		row := displayRowFuncs[key]

		data := endpoint.NewEndpointData[*cost.Cost](store, nil, filterFuncs)
		display := endpoint.NewEndpointDisplay[*cost.Cost](head, row, nil)
		ep := endpoint.New[*cost.Cost]("aws-cost-monthly-monthly-total", response, data, display, parameters)

		server.Middleware(ep.ProcessRequest, server.LoggingMW, server.SecurityHeadersMW)(w, r)
	})

	// Previously "Service" sheet
	//
	// 	- /aws/costs/{version}/monthly/{start}/{end}/units/{$}
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/{$}", func(w http.ResponseWriter, r *http.Request) {
		response := resp.New()
		key := "perUnit"

		parameters := qp.Parse(r)
		filterFuncs := FilterFunctions(parameters, response)
		displayHeadFuncs := DisplayHeadFunctions(parameters)
		displayRowFuncs := DisplayRowFunctions(parameters)
		displayFootFuncs := DisplayFootFunctions(parameters)

		head := displayHeadFuncs[key]
		row := displayRowFuncs[key]
		foot := displayFootFuncs[key]

		data := endpoint.NewEndpointData[*cost.Cost](store, byUnit, filterFuncs)
		display := endpoint.NewEndpointDisplay[*cost.Cost](head, row, foot)
		ep := endpoint.New[*cost.Cost]("aws-cost-monthly-per-unit", response, data, display, parameters)

		server.Middleware(ep.ProcessRequest, server.LoggingMW, server.SecurityHeadersMW)(w, r)
	})

	// Previously "Service And Environment" sheet
	//
	// 	- /aws/costs/{version}/monthly/{start}/{end}/units/envs/{$}
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/{$}", func(w http.ResponseWriter, r *http.Request) {
		response := resp.New()
		key := "perUnitEnv"

		parameters := qp.Parse(r)
		filterFuncs := FilterFunctions(parameters, response)
		displayHeadFuncs := DisplayHeadFunctions(parameters)
		displayRowFuncs := DisplayRowFunctions(parameters)
		displayFootFuncs := DisplayFootFunctions(parameters)

		head := displayHeadFuncs[key]
		row := displayRowFuncs[key]
		foot := displayFootFuncs[key]

		data := endpoint.NewEndpointData[*cost.Cost](store, byUnitEnv, filterFuncs)
		display := endpoint.NewEndpointDisplay[*cost.Cost](head, row, foot)
		ep := endpoint.New[*cost.Cost]("aws-cost-monthly-per-unit-env", response, data, display, parameters)

		server.Middleware(ep.ProcessRequest, server.LoggingMW, server.SecurityHeadersMW)(w, r)

	})

	// Previously "Detailed breakdown" sheet
	//
	// 	- /aws/costs/{version}/monthly/{start}/{end}/units/envs/services/{$}
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/services/{$}", func(w http.ResponseWriter, r *http.Request) {
		response := resp.New()
		key := "perUnitEnvService"

		parameters := qp.Parse(r)
		filterFuncs := FilterFunctions(parameters, response)
		displayHeadFuncs := DisplayHeadFunctions(parameters)
		displayRowFuncs := DisplayRowFunctions(parameters)
		displayFootFuncs := DisplayFootFunctions(parameters)

		head := displayHeadFuncs[key]
		row := displayRowFuncs[key]
		foot := displayFootFuncs[key]

		data := endpoint.NewEndpointData[*cost.Cost](store, byUnitEnvService, filterFuncs)
		display := endpoint.NewEndpointDisplay[*cost.Cost](head, row, foot)
		ep := endpoint.New[*cost.Cost]("aws-cost-monthly-per-unit-env-service", response, data, display, parameters)

		server.Middleware(ep.ProcessRequest, server.LoggingMW, server.SecurityHeadersMW)(w, r)
	})
}
