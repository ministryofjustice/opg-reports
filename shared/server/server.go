package server

import (
	"encoding/json"
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/server/response"
)

type IApiWithDataStore[V data.IEntry] interface {
	// Store returns the configured data store
	Store() data.IStore[V]
}
type IApiWithFileSystem[F files.IReadFS] interface {
	// FS will return the configure filesystem
	FS() F
}

type IApiRouteRegistration interface {
	// Register the routes this api handles to the mux
	Register(mux *http.ServeMux)
}

type IApiResponse[C response.ICell, R response.IRow[C]] interface {
	GetResponse() response.IResponse[C, R]
	GetNewResponse() response.IResponse[C, R]
}

type IApiWrite interface {
	Write(w http.ResponseWriter, status int, content []byte)
}

type IApiBase[V data.IEntry, F files.IReadFS] interface {
	IApiWithDataStore[V]
	IApiWithFileSystem[F]
	IApiRouteRegistration
	IApiWrite
}

// IApi is common interface for setting up a standard data driven api
type IApi[V data.IEntry, F files.IReadFS, C response.ICell, R response.IRow[C]] interface {
	IApiBase[V, F]
	IApiResponse[C, R]
	Start(w http.ResponseWriter, r *http.Request)
	End(w http.ResponseWriter, r *http.Request)
}

// --- CONCREATE

type Api[V data.IEntry, F files.IReadFS, C response.ICell, R response.IRow[C]] struct {
	store data.IStore[V]
	fs    F
	resp  response.IResponse[C, R]
}

func (a *Api[V, F, C, R]) Store() data.IStore[V] {
	return a.store
}
func (a *Api[V, F, C, R]) FS() F {
	return a.fs
}

func (a *Api[V, F, C, R]) Register(mux *http.ServeMux) {}
func (a *Api[V, F, C, R]) Write(w http.ResponseWriter, status int, content []byte) {
	w.WriteHeader(status)
	w.Write(content)
}
func (a *Api[V, F, C, R]) GetResponse() response.IResponse[C, R] {
	return a.resp
}
func (a *Api[V, F, C, R]) GetNewResponse() response.IResponse[C, R] {
	a.resp = response.NewResponse[C, R]()
	return a.resp
}

func (a *Api[V, F, C, R]) Start(w http.ResponseWriter, r *http.Request) {
	a.GetNewResponse().SetStart()
}

func (a *Api[V, F, C, R]) End(w http.ResponseWriter, r *http.Request) {
	resp := a.GetResponse()
	resp.SetEnd()
	resp.SetDuration()
	resp.GetDataAgeMin()
	resp.GetDataAgeMax()

	content, _ := json.MarshalIndent(resp, "", "  ")
	a.Write(w, resp.GetStatus(), content)
}

func NewApi[V data.IEntry, F files.IReadFS, C response.ICell, R response.IRow[C]](
	store data.IStore[V],
	fileSys F,
	resp response.IResponse[C, R]) *Api[V, F, C, R] {

	return &Api[V, F, C, R]{
		store: store,
		fs:    fileSys,
		resp:  resp,
	}
}
