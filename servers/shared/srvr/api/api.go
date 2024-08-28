package api

import (
	"context"
	"net/http"
)

type ApiServer struct {
	Ctx    context.Context
	DbPath string
}

type ApiHandlerFunc func(server *ApiServer, w http.ResponseWriter, r *http.Request)

// Wrap wraps a known function in a HandlerFunc and passes along the server details it needs
// outside of the normal scope of the http request
// Ths resulting function is then passed to the mux handle func (or via middleware)
func Wrap(server *ApiServer, innerFunc ApiHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		innerFunc(server, w, r)
	}
}

func New(ctx context.Context, dbPath string) *ApiServer {
	return &ApiServer{
		Ctx:    ctx,
		DbPath: dbPath,
	}
}
