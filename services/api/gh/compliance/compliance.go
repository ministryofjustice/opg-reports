package compliance

import (
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/gh/comp"
	"opg-reports/shared/server"
)

// Api is a concreate version
type Api[V *comp.Compliance, F files.WriteFS] struct {
	store *data.Store[*comp.Compliance]
	fs    *files.WriteFS
}

func (a *Api[V, F]) Store() data.IStore[*comp.Compliance] {
	return a.store
}

func (a *Api[V, F]) FS() files.IWriteFS {
	return a.fs
}

func (a *Api[V, F]) Register(mux *http.ServeMux) {

	mux.HandleFunc("/github/compliance/{version}/list/{$}",
		server.Middleware(a.List, server.LoggingMW, server.SecurityHeadersMW))
}

func (a *Api[V, F]) Write(w http.ResponseWriter, status int, content []byte) {
	w.WriteHeader(status)
	w.Write(content)
}

func New[V *comp.Compliance, F files.WriteFS](store *data.Store[*comp.Compliance], fS *files.WriteFS) *Api[V, F] {

	return &Api[V, F]{
		store: store,
		fs:    fS,
	}

}
