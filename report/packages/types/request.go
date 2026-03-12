package types

import "net/http"

// Parameters is the underlying fields that get populated
// via processing the incoming http request.
//
// Use in converting http request into a sql filter and
// is generally attached to the result
type Parameters interface {
	// Keys returns all the json names of each field on the
	// struct as slice so that can be iterated over / checked
	Keys() []string
	// Data returns a map of the struct data, uses json
	// marshaling so the keys match the `json` tag.
	//
	// Used to help create a filter and attached to the
	// server response
	Data() map[string]string
}

// HttpRequest exposes methods to set / get the original
// http request struct
type HttpRequester interface {
	// SetRequest attaches the http request to the
	// struct which can also then update itself and
	// configure any values
	SetRequest(r *http.Request)
	// Request returns the original http request
	Request() *http.Request
}

// ParameterRequester exposes methods to fetch the
// parameter data that was generated from the
// original http.request
type ParameterRequester interface {
	// Parameters returns the allowed and processed
	// data generated from the http request
	Parameters() Parameters
}

// ServerRequester is base line request passed along inside
// the mux handler.
type ServerRequester interface {
	HttpRequester
	ParameterRequester
}
