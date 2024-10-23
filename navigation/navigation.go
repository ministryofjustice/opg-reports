package navigation

import "github.com/ministryofjustice/opg-reports/endpoints"

type NavigationParameterFunc func(uri string, args ...interface{}) string

// NavigationDisplay contains all information relating to display
// of a Navigation struct
type NavigationDisplay struct {
	IsHeader   bool `json:"is_header" doc:"When true, the this is treated as a header item and not a standard link."`
	IsActive   bool `json:"is_active" doc:"Only true when this page is an exact match for the current uri on front."`
	IsInUri    bool `json:"is_in_uri" doc:"Will be true when the uri attached is with the url path (so /test/ is in /test/ab/123)."`
	Registered bool `json:"registered" doc:"Used to mark if the navigation struct has been registered by the handler"`
}

type NavigationData struct {
	Source  endpoints.ApiEndpoint
	Parsers []NavigationParameterFunc
}

// Navigation represents a tree structured navigation that is rendered on the
// front end of website to show sections at a time
// It also contains details on how to fetch data for each of the pages it represents
type Navigation struct {
	Name     string             `json:"name" doc:"Name is used when displaying this navigation link."`
	Uri      string             `json:"uri" doc:"Uri is front end url this navigation item will render for."`
	Display  *NavigationDisplay `json:"display" doc:"rendering related details"`
	Data     []*NavigationData  `json:"data" doc:"list of api endpoints to get data from"`
	Children []*Navigation      `json:"children" doc:"Child navigation"`
}
