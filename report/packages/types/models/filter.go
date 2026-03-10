package models

import "opg-reports/report/packages/convert"

// Filter is used for the api handler for any sql
// and contains all possible filters options.
//
// Empty fields are ignored.
type Filter struct {
	Team   string   `json:"team,omitempty"`   // Team filter
	Months []string `json:"months,omitempty"` // Months is generated from DateStart..DateEnd or DateA + DateB
	// Archived string   `json:"archived,omitempty"` // Archived bool - when "true" shows only archived code
}

// Map returns the values of field as a map; used for select
func (self *Filter) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}
