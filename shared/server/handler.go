package server

import "net/http"

// HttpHandler matches func signature for http.ServeMux.HandleFunc handler
type HttpHandler func(w http.ResponseWriter, r *http.Request)
