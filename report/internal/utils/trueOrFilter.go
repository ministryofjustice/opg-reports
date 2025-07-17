package utils

import "strings"

// TrueOrFilter is used to allow a query string parameter
// to be set as true or a value and how it should be used
// within sql statement
//
// When true, should be used in the selection, ordering &
// grouping of data, as in this column is enabled
//
// When the value is something other than "" or "true" then
// the field should be treat as a filter, being used
// in the select & where clause only
//
// `?team=true` => group the data by "team"
// `?team=my-team` => filter the data by team='my-team'
//
// Reduces the number of statements and query params
// needed.
type TrueOrFilter string

// Selectable returns true as long as the field is
// not empty
func (self TrueOrFilter) Selectable() bool {
	return (self != "")
}

// Groupable only returns true when the field value
// is a string "true"
func (self TrueOrFilter) Groupable() bool {
	return (strings.ToLower(string(self)) == "true")
}

// Orderable only returns true when the field value
// is a string "true"
func (self TrueOrFilter) Orderable() bool {
	return (strings.ToLower(string(self)) == "true")
}

// Whereable is only true when value is a non-empty string
// that is not "true"
func (self TrueOrFilter) Whereable() bool {
	return (self != "" && strings.ToLower(string(self)) != "true")
}
