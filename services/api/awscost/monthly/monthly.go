package monthly

import (
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/server"
)

// Api is a concreate version
type Api[V data.IEntry, F files.IReadFS] struct {
	store    data.IStore[V]
	fs       F
	response *ApiResponse
}

func (a *Api[V, F]) Store() data.IStore[V] {
	return a.store
}

func (a *Api[V, F]) FS() F {
	return a.fs
}
func (a *Api[V, F]) Response() server.IApiResponse {
	return a.response
}

func (a *Api[V, F]) Register(mux *http.ServeMux) {
	mux.HandleFunc("/aws/costs/{version}/monthly", a.Index)
}

func (a *Api[V, F]) Write(w http.ResponseWriter, response server.IApiResponse) {
	w.WriteHeader(response.Status())
	w.Write(response.Body())
}

func (a *Api[V, F]) Index(w http.ResponseWriter, r *http.Request) {
	res := a.Response()
	res.Start()

	res.Set(a.store.List(), http.StatusOK, nil)

	res.End()
	a.Write(w, res)
}

func New[V data.IEntry, D data.IStore[V], F files.IReadFS](store D, fS F) *Api[V, F] {
	return &Api[V, F]{
		store:    store,
		fs:       fS,
		response: &ApiResponse{},
	}

}
