package interfaces

import (
	"net/http"
	"opg-reports/report/packages/args"
)

// Result is end result returned from the api
type Result map[string]interface{}

type Resultable interface {
	Result() Result
}

// Row used to handle the results from the database
// select statements.
type Row interface {
	Selectable
	Resultable
}

// Resetter is used to reset response values on the struct
type Resetter interface {
	Reset()
}

// HttpWriter exposes function to get the http response writer interface
type HttpWriter interface {
	// Writer sets the response
	Writer(w http.ResponseWriter)
	// Response returns the http response writer interface
	Response() http.ResponseWriter
}

// Versioner exposes function to set version information on the struct
type Versioner interface {
	// Version sets the version information
	Version(version *args.Versions)
}

type Typed interface {
	Typed(t string)
}

// Requester exposes function to set the incoming request to the response
type Requester interface {
	// Request sets the original incoming reqest details
	Request(request ApiRequest)
}

// Filterer exposes function to set the filter details to the resposne
type Filterer interface {
	// Filter sets the filter details to the resposne
	Filter(filter Filterable)
}

type Databaser interface {
	// Record appends the db result
	Record(record Row)
	// Records returns results that have been set
	Records() []Row
}

type Resulter interface {
	Result(result Result)
	Results() []Result
}

// ApiResponse is the main interface used to respond to
// api requests and return data
type ApiResponse interface {
	Resetter
	HttpWriter
	Versioner
	Typed
	Requester
	Filterer
	Databaser
	Resulter
}
