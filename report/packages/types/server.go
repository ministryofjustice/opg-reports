package types

import "net/http"

// Configurable exposes methods to updte and retrieve the servers
// configuration details
type Configurable interface {
	SetConfig(c ServerConfigure)
	Config() ServerConfigure
}

// Muxable allows setting of the server mux struct
// onto this struct to be able to then attach endpoints
// and so on
type Muxable interface {
	SetMux(mux *http.ServeMux)
	Mux() *http.ServeMux
}
