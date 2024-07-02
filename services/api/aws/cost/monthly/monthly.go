package monthly

import (
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/server"
	"opg-reports/shared/server/response"
)

type MonthlyApi[V data.IEntry, F files.IReadFS, C response.ICell, R response.IRow[C]] struct {
	server.IApi[V, F, C, R]
}

func (a *MonthlyApi[V, F, C, R]) Register(mux *http.ServeMux) {
	// mux.HandleFunc("/aws/costs/{version}/monthly/{$}",
	// 	server.Middleware(a.Index, server.LoggingMW, server.SecurityHeadersMW))
	// mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/{$}",
	// 	server.Middleware(a.Totals, server.LoggingMW, server.SecurityHeadersMW))
	// mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/{$}",
	// 	server.Middleware(a.Units, server.LoggingMW, server.SecurityHeadersMW))
	// mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/{$}",
	// 	server.Middleware(a.UnitEnvironments, server.LoggingMW, server.SecurityHeadersMW))
	// mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/units/envs/services/{$}",
	// 	server.Middleware(a.UnitEnvironmentServices, server.LoggingMW, server.SecurityHeadersMW))
}
func NewMonthlyApi[V data.IEntry, F files.IReadFS, C response.ICell, R response.IRow[C]](
	store data.IStore[V],
	fileSys F,
	resp response.IResponse[C, R]) *MonthlyApi[V, F, C, R] {

	api := server.NewApi[V, F, C, R](store, fileSys, resp)
	return &MonthlyApi[V, F, C, R]{IApi: api}

}

// Api is a concreate version
type Api[V *cost.Cost, F files.WriteFS] struct {
	store *data.Store[*cost.Cost]
	fs    *files.WriteFS
}

func (a *Api[V, F]) Store() data.IStore[*cost.Cost] {
	return a.store
}

func (a *Api[V, F]) FS() files.IWriteFS {
	return a.fs
}

func (a *Api[V, F]) Register(mux *http.ServeMux) {
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

func (a *Api[V, F]) Write(w http.ResponseWriter, status int, content []byte) {
	w.WriteHeader(status)
	w.Write(content)
}

func New[V *cost.Cost, F files.WriteFS](store *data.Store[*cost.Cost], fS *files.WriteFS) *Api[V, F] {

	return &Api[V, F]{
		store: store,
		fs:    fS,
	}

}
