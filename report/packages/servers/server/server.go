package server

import (
	"net/http"
	"opg-reports/report/packages/types"
)

type server struct {
	// ctx is the context & logger that this struct
	// utilises.
	//
	// Set only during creation.
	ctx types.ContextLogger
	// config used to store details about this
	// server setup such as the hostname and
	// database configuration
	//
	// Set only during creation.
	config types.ServerConfigure
	// mux is the generated mux to use for all
	// handler mapping
	//
	// Set only during creation.
	mux *http.ServeMux
}

// SetCtx
func (self *server) SetCtx(ctx types.ContextLogger) {
	self.ctx = ctx
}

// Ctx
func (self *server) Ctx() types.ContextLogger {
	return self.ctx
}

// SetConfig
func (self *server) SetConfig(cfg types.ServerConfigure) {
	self.config = cfg
}

// Config
func (self *server) Config() types.ServerConfigure {
	return self.config
}

// SetMux
func (self *server) SetMux(mux *http.ServeMux) {
	self.mux = mux
}

// Mux
func (self *server) Mux() *http.ServeMux {
	return self.mux
}
