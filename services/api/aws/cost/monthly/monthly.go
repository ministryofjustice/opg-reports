package monthly

import (
	"log/slog"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/files"
	"opg-reports/shared/server"
	"opg-reports/shared/server/endpoint"
	"opg-reports/shared/server/resp"
	"opg-reports/shared/server/response"
	"time"
)

type Api[V *cost.Cost, F files.IReadFS, C response.ICell, R response.IRow[C]] struct {
	*server.Api[*cost.Cost, F, C, R]
}

func (a *Api[V, F, C, R]) Register(mux *http.ServeMux) {
	mux.HandleFunc("/aws/costs/{version}/monthly/{$}",
		server.Middleware(a.Index, server.LoggingMW, server.SecurityHeadersMW))
	mux.HandleFunc("/aws/costs/v1/monthly/{start}/{end}/{$}",
		server.Middleware(a.MonthlyTotals, server.LoggingMW, server.SecurityHeadersMW))
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/{$}",
		server.Middleware(a.MonthlyCostsPerAccountUnits, server.LoggingMW, server.SecurityHeadersMW))
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/{$}",
		server.Middleware(a.MonthlyCostsPerAccountUnitAndEnvironments, server.LoggingMW, server.SecurityHeadersMW))
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/services/{$}",
		server.Middleware(a.MonthlyCostsPerAccountUnitEnvironmentAndServices, server.LoggingMW, server.SecurityHeadersMW))
}

func (a *Api[V, F, C, R]) startAndEndDates(r *http.Request) (startDate time.Time, endDate time.Time) {
	var err error
	res := a.GetResponse()
	now := time.Now().UTC().Format(dates.FormatYM)
	// Get the dates
	startDate, err = dates.StringToDateDefault(r.PathValue("start"), "-", now)
	if err != nil {
		res.SetErrorAndStatus(err, http.StatusConflict)
	}
	endDate, err = dates.StringToDateDefault(r.PathValue("end"), "-", now)
	if err != nil {
		res.SetErrorAndStatus(err, http.StatusConflict)
	}
	slog.Info("[api/aws/costs/monthly] start and end dates",
		slog.Time("start", startDate),
		slog.Time("end", endDate))
	return

}

func New[V *cost.Cost, F files.IReadFS, C response.ICell, R response.IRow[C]](
	store data.IStore[*cost.Cost],
	fileSys F,
	resp response.IResponse[C, R]) *Api[*cost.Cost, F, C, R] {

	api := server.NewApi[*cost.Cost, F, C, R](store, fileSys, resp)
	return &Api[*cost.Cost, F, C, R]{Api: api}

}

const taxServiceName string = "tax"

var allowedParameters = []string{
	"start",
	"end",
	"unit",
	"environment",
}

func Register(mux *http.ServeMux, store data.IStore[*cost.Cost]) {

	qp := endpoint.NewQueryable(allowedParameters)
	// MonthlyTotals: /aws/costs/{version}/monthly/{start}/{end}/{$}
	// Returns cost data split into with & without tax segments, then grouped by the month
	// Previously "Totals" sheet
	// Note: if {start} or {end} are "-" it uses current month
	mux.HandleFunc("/aws/costs/v2/monthly/{start}/{end}/{$}", func(w http.ResponseWriter, r *http.Request) {
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
		ep := endpoint.New[*cost.Cost]("test", response, data, display, parameters)

		server.Middleware(ep.ProcessRequest, server.LoggingMW, server.SecurityHeadersMW)(w, r)
	})

}
