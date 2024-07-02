package monthly

import (
	"log/slog"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/files"
	"opg-reports/shared/server"
	"opg-reports/shared/server/response"
	"time"
)

type Api[V *cost.Cost, F files.IReadFS, C response.ICell, R response.IRow[C]] struct {
	*server.Api[*cost.Cost, F, C, R]
}

func (a *Api[V, F, C, R]) Register(mux *http.ServeMux) {
	mux.HandleFunc("/aws/costs/{version}/monthly/{$}",
		server.Middleware(a.Index, server.LoggingMW, server.SecurityHeadersMW))
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/{$}",
		server.Middleware(a.Totals, server.LoggingMW, server.SecurityHeadersMW))
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/{$}",
		server.Middleware(a.Units, server.LoggingMW, server.SecurityHeadersMW))
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/{$}",
		server.Middleware(a.UnitEnvironments, server.LoggingMW, server.SecurityHeadersMW))
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/services/{$}",
		server.Middleware(a.UnitEnvironmentServices, server.LoggingMW, server.SecurityHeadersMW))
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

// // Api is a concreate version
// type Api[V *cost.Cost, F files.WriteFS] struct {
// 	store *data.Store[*cost.Cost]
// 	fs    *files.WriteFS
// }

// func (a *Api[V, F]) Store() data.IStore[*cost.Cost] {
// 	return a.store
// }

// func (a *Api[V, F]) FS() files.IWriteFS {
// 	return a.fs
// }

// func (a *Api[V, F]) Register(mux *http.ServeMux) {
// 	mux.HandleFunc("/aws/costs/{version}/monthly/{$}",
// 		server.Middleware(a.Index, server.LoggingMW, server.SecurityHeadersMW))
// 	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/{$}",
// 		server.Middleware(a.Totals, server.LoggingMW, server.SecurityHeadersMW))
// 	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/{$}",
// 		server.Middleware(a.Units, server.LoggingMW, server.SecurityHeadersMW))
// 	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/{$}",
// 		server.Middleware(a.UnitEnvironments, server.LoggingMW, server.SecurityHeadersMW))
// 	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/services/{$}",
// 		server.Middleware(a.UnitEnvironmentServices, server.LoggingMW, server.SecurityHeadersMW))
// }

// func (a *Api[V, F]) Write(w http.ResponseWriter, status int, content []byte) {
// 	w.WriteHeader(status)
// 	w.Write(content)
// }

// func New[V *cost.Cost, F files.WriteFS](store *data.Store[*cost.Cost], fS *files.WriteFS) *Api[V, F] {

// 	return &Api[V, F]{
// 		store: store,
// 		fs:    fS,
// 	}

// }
