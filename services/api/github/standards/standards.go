package standards

import (
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/github/std"
	"opg-reports/shared/server"
	"opg-reports/shared/server/response"
)

// Api implements IApi
type Api[V *std.Repository, F files.IReadFS, C response.ICell, R response.IRow[C]] struct {
	*server.Api[*std.Repository, F, C, R]
}

func (a *Api[V, F, C, R]) Register(mux *http.ServeMux) {
	mux.HandleFunc("/github/standards/{version}/list/{$}",
		server.Middleware(a.List, server.LoggingMW, server.SecurityHeadersMW))
}

// AllowedGetParameters allows this data to be filtered by
// - archived
// - teams
func (a *Api[V, F, C, R]) AllowedGetParameters() []string {
	return []string{
		"archived",
		"teams",
	}
}

func New[V *std.Repository, F files.IReadFS, C response.ICell, R response.IRow[C]](
	store data.IStore[*std.Repository],
	fileSys F,
	resp response.IResponse[C, R]) *Api[*std.Repository, F, C, R] {

	api := server.NewApi[*std.Repository, F, C, R](store, fileSys, resp)
	return &Api[*std.Repository, F, C, R]{Api: api}

}
