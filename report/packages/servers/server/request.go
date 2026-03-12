package server

import (
	"net/http"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types"
	"reflect"
	"strings"
)

// parameters struct contains the list of automatically loaded values from
// the http.request
type parameters struct {
	DateStart string `json:"date_start,omitempty"` // start date for any date ranges
	DateEnd   string `json:"date_end,omitempty"`   // end date for and date ranges
	DateA     string `json:"date_a,omitempty"`     // DateA used in date comparisons
	DateB     string `json:"date_b,omitempty"`     // DateB used in date comparisons
	Team      string `json:"team,omitempty"`       // Team name filters
}

// Data returns the values of this struct as a map using
// json marshaling
func (self *parameters) Data() (values map[string]string) {
	values = map[string]string{}
	convert.Between(self, &values)
	return
}

// Keys returns the name of each field on the struct
// using reflection
func (self *parameters) Keys() (keys []string) {
	var t = reflect.TypeOf(self).Elem()
	keys = []string{}

	for i := 0; i < t.NumField(); i++ {
		var tag = t.Field(i).Tag.Get("json")
		var name = strings.Split(tag, ",")[0]
		if name != "-" {
			keys = append(keys, name)
		}
	}
	return
}

// Request is a struct to handle processing http.request
// into known allowed fields.
//
// Looks are a list of allowed parameter names in both the
// path and query string and when found sets the value to
// match. Fetch this via `.Parameters().Data()`
//
// Provides methods to allow accessing the original http
// request. Satisfies the following interfaces:
//
// types.HttpRequester
// types.ParameterRequester
// types.ServerRequester
type Request struct {
	// req is internal attribute to track the original http request
	// struct that is passed from mux handlers etc
	req *http.Request `json:"-"`
	// parameters
	parameters *parameters `json:"-"`
}

// update is used fetch values from the http.request into
// this struct so they can used by the server handler.
//
// Replaces the parameters struct with an updated version
// that is populated from the http.request values.
//
// For each field it will look for the path value first and then
// attempt the query string (returning first value only).
func (self *Request) update(r *http.Request) {
	var params = &parameters{}
	var data = map[string]string{}
	var query = r.URL.Query()

	for _, key := range params.Keys() {
		if val := r.PathValue(key); val != "" {
			data[key] = val
		} else if val := query.Get(key); val != "" {
			data[key] = val
		}
	}
	// convert from the map to the params
	convert.Between(data, &params)
	// update self
	self.parameters = params
}

// Parameters returns the underlying data processed from
// the http request
func (self *Request) Parameters() types.Parameters {
	return self.parameters
}

// SetRequest attach the request and update parameter values
func (self *Request) SetRequest(r *http.Request) {
	self.req = r
	self.update(r)
}

// Request returns the original http request
func (self *Request) Request() *http.Request {
	return self.req
}
