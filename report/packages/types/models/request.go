package models

import (
	"net/http"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/ranges"
	"opg-reports/report/packages/types/interfaces"
)

// Request used with the api handlers and wraps http.Request.
//
// Contains all of the incoming options that api handler
// endpoint could utilise.
//
// Empty fields are ignored.
//
// interfaces.HttpRequester
// interfaces.Populator
// interface.FilterMaker
type Request struct {
	Req *http.Request `json:"-"`   // used for the request getter
	Uri string        `json:"uri"` // the original query string called

	DateStart string `json:"date_start"` // start date for any date ranges
	DateEnd   string `json:"date_end"`   // end date for and date ranges
	DateA     string `json:"date_a"`     // DateA used in date comparisons
	DateB     string `json:"date_b"`     // DateB used in date comparisons
	Team      string `json:"team"`       // Team name filters
	// Archived  string `json:"archived"`   // Archived, when set, will hide archived code bases
}

// fields returns the json field names for the request struct
func (self *Request) fields() (fields []string) {
	var m = map[string]interface{}{}
	convert.Between(self, &m)
	fields = []string{}
	for k, _ := range m {
		fields = append(fields, k)
	}
	return
}

// HttpRequest sets the request
func (self *Request) HttpRequest(r *http.Request) {
	self.Req = r
	self.Uri = r.RequestURI
}
func (self *Request) Request() *http.Request {
	return self.Req
}

// Populate iterates over its own fields and looks for those values in the
// request object, first from the path then the query string.
//
// Uses `url.Values.Get` to look in query string data, so only returns
// the first value rather than a slice of all
func (self *Request) Populate(req *http.Request) {
	var data = map[string]interface{}{}
	var fields = self.fields()
	var query = req.URL.Query()

	for _, field := range fields {
		// try the path value first, then check the query
		// string values - only fetches the first
		if val := req.PathValue(field); val != "" {
			data[field] = val
		} else if qs := query.Get(field); qs != "" {
			data[field] = qs
		}
	}
	// update self
	convert.Between(data, &self)

}

// Filter creates a interfaces.Filterable struct combination of itself and
// the incoming http request.
//
// Will convert date_start & date_end into Months and similary will
// set date_a & date_b as Months as well.
//
// DateStart & DateEnd take precidence over DateA & DateB.
func (self *Request) Filter(req *http.Request) interfaces.Filterable {
	var filter = &Filter{}
	// do the convert to handle matching json keys
	convert.Between(self, &filter)
	// now setup the more elaborate values on the filter
	// - months from the start & end
	// - months from date a & b
	if self.DateStart != "" && self.DateEnd != "" {
		filter.Months = ranges.Months(self.DateStart, self.DateEnd)
	} else if self.DateA != "" && self.DateB != "" {
		filter.Months = []string{
			self.DateA, self.DateB,
		}
	}

	return filter
}
