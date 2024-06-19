package server

import (
	"net/http"
	"opg-reports/shared/files"
)

type IServer[F files.IReadFS] interface {
	// FS will return the configure filesystem
	FS() F
	//
	Response() IResponse
	// Register the routes this api handles to the mux
	Register(mux *http.ServeMux)
	// Write outputs status and data
	Write(w http.ResponseWriter, response IResponse)
}
