package httpx

import (
	"net/http"
	"opg-reports/report/packages/types"
)

// Request extends http.Request with new methods to help process
// server calls
//
// TODO:
//   - ADD CTX & LOGGING
type Request struct {
	*http.Request
}

// HttpRequest returns the original
func (self *Request) HttpRequest() *http.Request {
	return self.Request
}

// Parameters gets data from the request.
//
// iterates over fields and looks for it in the original
// request path or the query strings, when found attaches
// to the data map returned
func (self *Request) Parameters(fields []string) (data map[string]string) {
	var r = self.Request
	var qs = r.URL.Query()

	data = map[string]string{}

	for _, field := range fields {
		if v := r.PathValue(field); v != "" {
			data[field] = v
		} else if v := qs.Get(field); v != "" {
			data[field] = v
		}

	}

	return
}

// NewRequest returns request wrapped in the extended type
func NewRequest(r *http.Request) types.HttpxRequest {
	return &Request{Request: r}
}
