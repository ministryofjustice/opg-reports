// Packacge httpx is an extension of http
package httpx

import (
	"net/http"
	"opg-reports/report/packages/types"
)

// ServeMux an extended version of http.ServeMux
// that passes more data into the handler.
//
// Adds ..
//   - context
//
// TODO:
//   - logging
type ServeMux struct {
	*http.ServeMux
	// ctx is an extended context
	ctx types.Contextx
}

// HandleFuncx extends on the http.HandlerFuc type to pass along the
// context and template values.
func (self *ServeMux) HandleFuncx(
	pattern string,
	template types.HttpxTemplater,
	handler types.HttpxHandlerFunc) {

	// call the original mux handler and then call the custom handler inside of that
	self.ServeMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		var request = NewRequest(r)
		var writer = NewResponseWriter(w)
		// call the custom handler with template values and the updated
		// versions of writer and request
		handler(
			self.ctx,
			template,
			writer,
			request,
		)
	})
}

// NewServeMux returns custom mux version with extended
// capabilities
func NewServeMux() (mux *ServeMux) {
	mux = &ServeMux{
		ServeMux: http.NewServeMux(),
	}
	return
}
