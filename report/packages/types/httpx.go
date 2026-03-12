package types

import (
	"html/template"
	"net/http"
)

// TODOS:
// 	 - BUILD REQUEST HANDLER (start with api)
// 	 - BUILD LOGGING into context
// 	 - BUILD CONFIG
// 		- build context func to get config from context

// HttpxRequest is an interface for a extended *http.Request.
//
// Adds function to pull parameters directly from the query
// so dont have to do that each time
//
// Puts some existing methods in http.Request to the interface
// that are used within Parameters
//
// extends: *http.Request
type HttpxRequest interface {
	PathValue(name string) string
	// HttpRequest returns the original request
	HttpRequest() *http.Request
	// Parameters is extra, added to allow
	Parameters(fields []string) (data map[string]string)
}

// HttpxResponseWriter is an extension of http.ResponseWriter
type HttpxResponseWriter interface {
	http.ResponseWriter
	// Send wraps around the base Write to either process the
	// response into json / html before calling write
	Send(data any, template HttpxTemplater) (code int, err error)
}

// HttpxTemplater exposes method to access a complied
// html template struct that will be used by the extended
// writer interface for easier swapping between html and
// json responses.
type HttpxTemplater interface {
	Template() *template.Template
	Name() string
	WithName(n string) HttpxTemplater
}

// HttpxHandlerFunc is an extended handler that more data will
// get passed into.
//
// HttpxHandlerFunc method should call w.Send to instead of w.Write.
//
// extends: http.HandlerFunc
type HttpxHandlerFunc func(ctx Contextx, tpl HttpxTemplater, w HttpxResponseWriter, r HttpxRequest)

// HttpxMux is the extended serveMux with extra capability
type HttpxMux interface {
	HandleFuncx(pattern string, tpl HttpxTemplater, handler HttpxHandlerFunc)
}
