package interfaces

import "net/http"

type HttpRequester interface {
	HttpRequest(r *http.Request)
	Request() *http.Request
}

// Populator is interface used to populate a struct from an
// incoming http request.
//
// Values should be mapped from the request path values and
// query string into the local struct.
//
// Used for api handling to create a struct which contains
// incoming values
type Populator interface {
	// Populate uses the http request to update itself with values
	// from the request that match the json field names of its
	// own properties.
	Populate(req *http.Request)
}

// FilterMaker interfacee is used to help convert from a http request
// into a struct to use for filtering in the database
type FilterMaker interface {
	// Filter creates a Filterable struct from the http request object
	// and its self
	Filter(req *http.Request) Filterable
}

type ApiRequest interface {
	HttpRequester
	Populator
	FilterMaker
}
