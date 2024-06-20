package monthly

import (
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/server"
)

// Api is a concreate version
type Api[V *cost.Cost, F files.WriteFS] struct {
	store    *data.Store[*cost.Cost]
	fs       *files.WriteFS
	response *ApiResponse
}

func (a *Api[V, F]) Store() data.IStore[*cost.Cost] { //*data.Store[*cost.Cost]
	return a.store
}

func (a *Api[V, F]) FS() files.IWriteFS {
	return a.fs
}
func (a *Api[V, F]) Response() server.IApiResponse {
	return a.response
}

func (a *Api[V, F]) Register(mux *http.ServeMux) {
	mux.HandleFunc("/aws/costs/{version}/monthly/{$}", a.Index)
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/{$}", a.Totals)
}

func (a *Api[V, F]) Write(w http.ResponseWriter, response server.IApiResponse) {
	w.WriteHeader(response.GetStatus())
	w.Write(response.Body())
}

func New[V *cost.Cost, F files.WriteFS](store *data.Store[*cost.Cost], fS *files.WriteFS) *Api[V, F] {
	return &Api[V, F]{
		store:    store,
		fs:       fS,
		response: &ApiResponse{Status: http.StatusOK},
	}

}
