package httpx

import (
	"fmt"
	"net/http"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/ranges"
	"strings"
)

var standardFields = []string{
	`date_start`,
	`date_end`,
	`data_a`,
	`data_b`,
	`team`,
}

// RequestData represents the allow request path values and query strings that
// all server endpoints will utilise and provides them as struct fields
//
// Tracks the original request
type RequestData struct {
	r         *http.Request `json:"-"`
	DateStart string        `json:"date_start,omitempty"` // start date for any date ranges
	DateEnd   string        `json:"date_end,omitempty"`   // end date for and date ranges
	DateA     string        `json:"date_a,omitempty"`     // DateA used in date comparisons
	DateB     string        `json:"date_b,omitempty"`     // DateB used in date comparisons
	Team      string        `json:"team,omitempty"`       // Team name filters
}

// QueryString
// Used in apiclient.QueryStringer
func (self *RequestData) QueryString() (s string) {
	s = ""
	for k, v := range self.Map() {
		s += fmt.Sprintf("%s=%s&", k, v)
	}
	if len(s) > 0 {
		s = "?" + s
	}
	s = strings.Trim(s, "&")
	return
}

// Request returns the original request
func (self *RequestData) Request() *http.Request {
	return self.r
}

func (self *RequestData) Map() (m map[string]string) {
	m = map[string]string{}
	convert.Between(self, &m)
	return
}

// ValuesFromRequest returns a map of values fetched
// that match the keys passed whose values are pulled
// from the request path values or query string.
//
// Used to fetch incoming data such as start dates
// or team names that are then used to filter results
//
// Checks path value first, if not found it will
// then look in the query string.
//
// Values pulled from query string only return the first
// value.
func ValuesFromRequest(request *http.Request) (rd *RequestData) {
	var qs = request.URL.Query()
	var data = map[string]string{}
	rd = &RequestData{r: request}

	for _, key := range standardFields {
		if v := request.PathValue(key); v != "" {
			data[key] = v
		} else if v := qs.Get(key); v != "" {
			data[key] = v
		}
	}
	convert.Between(data, &rd)

	return
}

// RequestDataToFilter creates struct combination of itself and
// the incoming http request.
//
// Will convert date_start & date_end into Months and similary will
// set date_a & date_b as Months as well.
//
// DateStart & DateEnd take precidence over DateA & DateB.
func RequestDataToFilter(rd *RequestData) *Filter {
	var filter = &Filter{rd: rd}
	// do the convert to handle matching json keys
	convert.Between(rd, &filter)
	// now setup the more elaborate values on the filter
	// - months from the start & end
	// - months from date a & b
	if rd.DateStart != "" && rd.DateEnd != "" {
		filter.Months = ranges.Months(rd.DateStart, rd.DateEnd)
	} else if rd.DateA != "" && rd.DateB != "" {
		filter.Months = []string{
			rd.DateA, rd.DateB,
		}
	}

	return filter
}
