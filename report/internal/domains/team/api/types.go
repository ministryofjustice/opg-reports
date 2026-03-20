package api

import "opg-reports/report/packages/httpx"

// Result is wrapper for the api data results
type Result struct {
	httpx.ResponseData
	Teams []string `json:"data"`
}

// Team is used for the api and import setup to
// contain team data
type Team struct {
	Name string `json:"name"`
}

// Sequence is used to return the columns in the order they are selected
func (self *Team) Sequence() []any {
	return []any{&self.Name}
}
