package httpx

import (
	"net/http"
	"opg-reports/report/packages/convert"
)

type FitleredRequest interface {
	// Request returns the original http request that started this process
	Request() *http.Request
	// RequestData returns the processed http request struct which contains
	// the allowed set of fields used in filtering and so on
	RequestData() *RequestData
	// Filter returns the processed request data, where things like date_start
	// and date_end have been converted in to a list of months between those
	// dates
	Filter() *Filter
	// Map returns the filtered values as a map; used in the dbx.Select
	Map() (m map[string]interface{})
}

// Filter is used for the api handler for any sql
// and contains all possible filters options.
//
// Empty fields are ignored.
type Filter struct {
	rd     *RequestData `json:"-"`
	Team   string       `json:"team,omitempty"`   // Team filter
	Months []string     `json:"months,omitempty"` // Months is generated from DateStart..DateEnd or DateA + DateB
}

// Map returns the values of field as a map; used for select
func (self *Filter) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}

// Request returns the original http request that started this process
func (self *Filter) Request() *http.Request {
	return self.rd.Request()
}

// RequestData returns the processed http request struct which contains
// the allowed set of fields used in filtering and so on
func (self *Filter) RequestData() *RequestData {
	return self.rd
}

// Filter returns the processed request data, where things like date_start
// and date_end have been converted in to a list of months between those
// dates
func (self *Filter) Filter() *Filter {
	return self
}
