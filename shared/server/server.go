package server

import (
	"html/template"
	"net/http"
	"opg-reports/shared/data"
	"opg-reports/shared/files"
)

// IApi is common interface for setting up a standard data driven api
type IApi[V data.IEntry, F files.IReadFS] interface {
	// Store returns the configured data store
	Store() data.IStore[V]
	// FS will return the configure filesystem
	FS() F
	// Register the routes this api handles to the mux
	Register(mux *http.ServeMux)
	//
	Write(w http.ResponseWriter, status int, content []byte)
}

type IServer interface {
	// Register the routes this api handles to the mux
	Register(mux *http.ServeMux)
	//
	Write(w http.ResponseWriter, status int, tmpl *template.Template, name string, data any)
}