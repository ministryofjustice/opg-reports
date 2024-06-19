package monthly

import (
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/server"
)

// Api is a concreate version
type Api struct {
	store    *data.Store[*cost.Cost]
	fs       *files.WriteFS
	response *ApiResponse
}

func (a *Api) Store() *data.Store[*cost.Cost] {
	return a.store
}

func (a *Api) FS() *files.WriteFS {
	return a.fs
}
func (a *Api) Response() server.IApiResponse {
	return a.response
}

func (a *Api) Register(mux *http.ServeMux) {
	mux.HandleFunc("/aws/costs/{version}/monthly/{$}", a.Index)
	mux.HandleFunc("/aws/costs/{version}/monthly/{start}/{end}/{$}", a.Totals)
}

func (a *Api) Write(w http.ResponseWriter, response server.IApiResponse) {
	w.WriteHeader(response.Status())
	w.Write(response.Body())
}

func New(store *data.Store[*cost.Cost], fS *files.WriteFS) *Api {
	return &Api{
		store:    store,
		fs:       fS,
		response: &ApiResponse{},
	}

}
