package models

import (
	"net/http"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/types/interfaces"
)

// ApiResponseReceived is the simplifed response version of an api result
// mapping vales to simple types rather than interfaces
type ApiResponseReceived struct {
	Versions *args.Versions           `json:"versions"`
	Request  map[string]interface{}   `json:"request"`
	Filters  map[string]interface{}   `json:"filter"`
	Data     []map[string]interface{} `json:"data"`
}

type ApiResponse struct {
	// ResponseWriter is the original Http Writer we are wrapping around
	// within the custom mux to add more capabilities.
	ResponseWriter http.ResponseWriter `json:"-"`
	// Version data added from inputs options / env variables to help track
	// current code version.
	// Set from with the
	Versions *args.Versions `json:"versions"`
	// ID is simple text to identidy what endpoint was called and the
	// handler that processed it
	Type string `json:"type"`
	// OriginalRequest contains the request details that originally
	// called the endpoint
	OriginalRequest interfaces.ApiRequest `json:"request"`
	// Filters is the filter data used within this endpoint request
	// to adjust the sql select statement. Typically things like
	// data ranges and team names.
	Filters interfaces.Filterable `json:"filters"`
	// DatabaseResults are the raw results from the database.
	DatabaseResults []interfaces.Row `json:"-"`
	// The processed data that will be returned within the request.
	Data []interfaces.Result `json:"data"`
}

// Reset the response object values.
func (self *ApiResponse) Reset() {
	self.Type = ""
	self.OriginalRequest = nil
	self.Filters = nil
	self.DatabaseResults = []interfaces.Row{}
	self.Data = []interfaces.Result{}
}

// Writer sets the http response
func (self *ApiResponse) Writer(w http.ResponseWriter) {
	self.ResponseWriter = w
}

// Response returns the based repsonse writer
func (self *ApiResponse) Response() http.ResponseWriter {
	return self.ResponseWriter
}

// Version sets the verion details on this response
func (self *ApiResponse) Version(version *args.Versions) {
	self.Versions = version
}

// Typed
func (self *ApiResponse) Typed(t string) {
	self.Type = t
}

// Request attaches the original request details to this response struct
func (self *ApiResponse) Request(request interfaces.ApiRequest) {
	self.OriginalRequest = request
}

// Filter attaches the filter details to this response struct
func (self *ApiResponse) Filter(filter interfaces.Filterable) {
	self.Filters = filter
}

func (self *ApiResponse) Record(record interfaces.Row) {
	self.DatabaseResults = append(self.DatabaseResults, record)
}

func (self *ApiResponse) Records() []interfaces.Row {
	return self.DatabaseResults
}

func (self *ApiResponse) Result(result interfaces.Result) {
	self.Data = append(self.Data, result)
}

func (self *ApiResponse) Results() []interfaces.Result {
	return self.Data
}
