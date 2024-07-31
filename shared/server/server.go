// Package server provides a series of interfaces and concrete structs for request handling
//
// Api interfaces (IApi*) determine the features and functions that we require from one of
// our api "segments" that handle incoming requests. Our APIs are split based on area of
// concern, so areas like costs, standards and similar have their own concrete struct
// that impliments IApi.
//
// The package includes a concrete API struct that provides most of interface functionality,
// but does leave the `Register` func empty as this will be very specifc. Ideally, the
// localised concrete should use this Api struct as a base and provide overwrites where
// required to reduce effort of duplication as the api grows.
//
// The package also provides middleware functionality and some default middleware items.
// The Middleware function provides the chaining and should be used within the Register
// function of IApi. Simple logging and security header middlewares included.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
	"opg-reports/shared/server/response"
)

// IApiWithDataStore handles providing a data store
type IApiWithDataStore[V data.IEntry] interface {
	// Store returns the configured data store
	Store() data.IStore[V]
}

// IApiWithFileSystem provides a way to get the filesystem
type IApiWithFileSystem[F files.IReadFS] interface {
	// FS will return the configure filesystem
	FS() F
}

// IApiRouteRegistration handles registering routes for an api
type IApiRouteRegistration interface {
	// Register the routes this api handles to the mux
	Register(mux *http.ServeMux)
}

// IApiResponse provides getters & setters for setting what data will
// be sent back
type IApiResponse[C response.ICell, R response.IRow[C]] interface {
	GetResponse() response.IResponse[C, R]
	GetNewResponse() response.IResponse[C, R]
}

// IApiWrite provides method to write the data back to the client
type IApiWrite interface {
	Write(w http.ResponseWriter, status int, content []byte)
}

// IApiBase is a simple version of api with out response data
type IApiBase[V data.IEntry, F files.IReadFS] interface {
	IApiWithDataStore[V]
	IApiWithFileSystem[F]
	IApiRouteRegistration
	IApiWrite
}

type IApiQueryable[V data.IEntry] interface {
	AllowedGetParameters() []string
	GetParameters(allowed []string, r *http.Request) map[string][]string
}

// IApi is common interface for setting up a standard data driven api
// Adds Start to track the begining of the request and End to write the completed data back out
type IApi[V data.IEntry, F files.IReadFS, C response.ICell, R response.IRow[C]] interface {
	IApiBase[V, F]
	IApiResponse[C, R]
	IApiQueryable[V]
	Start(w http.ResponseWriter, r *http.Request)
	End(w http.ResponseWriter, r *http.Request)
}

// --- CONCREATE
// Api is a concrete api handler that provides typical functions to handle processing
// an incoming request, tracking its duration and provide methods to register handles and
// then send data back
// Interface: [IApi]
type Api[V data.IEntry, F files.IReadFS, C response.ICell, R response.IRow[C]] struct {
	store data.IStore[V]
	fs    F
	resp  response.IResponse[C, R]
}

// Store provides a data store to fetch info for processing within the request.
// Currently, this would be data read in from a series of files
func (a *Api[V, F, C, R]) Store() data.IStore[V] {
	slog.Debug("getting store")
	return a.store
}

// FS returns the reader filesystem helper for this running api - so can
// get any files etc
func (a *Api[V, F, C, R]) FS() F {
	slog.Debug("getting FS")
	return a.fs
}

// Register attached the relevent paths and handlers and registers their handlers
// to the mux passed
// Typically, this function will use middleware chainign etc
//
// Note: This version is empty and does nothing
func (a *Api[V, F, C, R]) Register(mux *http.ServeMux) {}

// Write uses the status and content to return info to the client
func (a *Api[V, F, C, R]) Write(w http.ResponseWriter, status int, content []byte) {
	slog.Debug("writing response", slog.Int("status", status))
	w.WriteHeader(status)
	w.Write(content)
}

// GetResponse returns the existing response item, allowing handlers to add result data
func (a *Api[V, F, C, R]) GetResponse() response.IResponse[C, R] {
	slog.Debug("return existing response")
	return a.resp
}

// GetNewResponse returns a brand new response - this should only be done within
// the Start() function, as it will blank any existing data that should be
// sent back in the response
func (a *Api[V, F, C, R]) GetNewResponse() response.IResponse[C, R] {
	slog.Debug("create new response")
	a.resp = response.NewResponse[C, R]()
	return a.resp
}

// AllowedGetParameters returns slice of strings of named GET parameters
// that can be used with this IApi.
// Defaults to empty
func (a *Api[V, F, C, R]) AllowedGetParameters() []string {
	return []string{}
}

// GetParameters uses the allowed slcie (the result of AllowedGetParameters) to find
// GET parameters that have been passed and returns a map of their values
// This is then used to build a filter
func (a *Api[V, F, C, R]) GetParameters(allowed []string, r *http.Request) map[string][]string {
	values := map[string][]string{}
	q := r.URL.Query()
	for _, name := range allowed {
		if v, ok := q[name]; ok {
			values[name] = v
		}
	}

	return values

}

// Start generates a new repsonse and then sets the start time of the request
func (a *Api[V, F, C, R]) Start(w http.ResponseWriter, r *http.Request) {
	slog.Debug("request start",
		slog.String("request_method", r.Method),
		slog.String("request_uri", r.URL.String()))

	a.GetNewResponse().SetStart()
}

// End fetches the response data, sets the end and duration of the request,
// works out the data age ranges for the information included an then
// writes the response out to as json to to the writer
func (a *Api[V, F, C, R]) End(w http.ResponseWriter, r *http.Request) {
	resp := a.GetResponse()
	resp.SetEnd()
	resp.SetDuration()
	resp.GetDataAgeMin()
	resp.GetDataAgeMax()

	content, _ := json.MarshalIndent(resp, "", "  ")
	status := resp.GetStatus()

	slog.Info("request end",
		slog.Int("status", status),
		slog.String("request_method", r.Method),
		slog.String("request_uri", r.URL.String()))

	a.Write(w, status, content)
}

// NewApi generates a new Api instance
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
