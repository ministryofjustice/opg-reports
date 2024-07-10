package standards

import (
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/github/std"
	"opg-reports/shared/server"
	"opg-reports/shared/server/response"
)

type Api[V *std.Repository, F files.IReadFS, C response.ICell, R response.IRow[C]] struct {
	*server.Api[*std.Repository, F, C, R]
}

func (a *Api[V, F, C, R]) Register(mux *http.ServeMux) {
	mux.HandleFunc("/github/standards/{version}/list/{$}",
		server.Middleware(a.List, server.LoggingMW, server.SecurityHeadersMW))
}

func New[V *std.Repository, F files.IReadFS, C response.ICell, R response.IRow[C]](
	store data.IStore[*std.Repository],
	fileSys F,
	resp response.IResponse[C, R]) *Api[*std.Repository, F, C, R] {

	api := server.NewApi[*std.Repository, F, C, R](store, fileSys, resp)
	return &Api[*std.Repository, F, C, R]{Api: api}

}

// // Api is a concreate version
// type Api[V *std.Repository, F files.WriteFS] struct {
// 	store *data.Store[*std.Repository]
// 	fs    *files.WriteFS
// }

// func (a *Api[V, F]) Store() data.IStore[*std.Repository] {
// 	return a.store
// }

// func (a *Api[V, F]) FS() files.IWriteFS {
// 	return a.fs
// }

// func (a *Api[V, F]) Register(mux *http.ServeMux) {

// }

// func (a *Api[V, F]) Write(w http.ResponseWriter, status int, content []byte) {
// 	w.WriteHeader(status)
// 	w.Write(content)
// }

// func New[V *std.Repository, F files.WriteFS](store *data.Store[*std.Repository], fS *files.WriteFS) *Api[V, F] {

// 	return &Api[V, F]{
// 		store: store,
// 		fs:    fS,
// 	}

// }
