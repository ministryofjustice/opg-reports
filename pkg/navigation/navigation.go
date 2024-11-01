// Package navigation used for navigation config on the front end.
//
// Allows for mapping a front end url to a series of data sources
// on the api and what template to render for each page on the
// front end.
package navigation

import (
	"github.com/ministryofjustice/opg-reports/pkg/endpoints"
)

// NavigationDisplay contains all information relating to display
// of a Navigation struct
type Display struct {
	IsHeader     bool   `json:"is_header" doc:"When true, the this is treated as a header item and not a standard link."`
	IsActive     bool   `json:"is_active" doc:"Only true when this page is an exact match for the current uri on front."`
	InUri        bool   `json:"in_uri" doc:"Will be true when the uri attached is with the url path (so /test/ is in /test/ab/123)."`
	Registered   bool   `json:"registered" doc:"Used to mark if the navigation struct has been registered by the handler"`
	PageTemplate string `json:"page_template" doc:"String name of the template to use for this page."`
}

type Data struct {
	Source    endpoints.ApiEndpoint `json:"source" faker:"uri" doc:"API data source location"`
	Namespace string                `json:"namepsace" faker:"word" doc:"Namespace the Result should be copied over into for front end parsing"`
	Body      interface{}           `json:"body" doc:"A pointer to the struct that this data would use as the body content."`
}

func NewData(source endpoints.ApiEndpoint, namespace string) *Data {
	return &Data{Source: source, Namespace: namespace}
}

// Navigation represents a tree structured navigation that is rendered on the
// front end of website to show sections at a time
// It also contains details on how to fetch data for each of the pages it represents
type Navigation struct {
	Name     string        `json:"name" doc:"Name is used when displaying this navigation link."`
	Uri      string        `json:"uri" doc:"Uri is front end url this navigation item will render for."`
	Display  *Display      `json:"display" doc:"rendering related details"`
	Data     []*Data       `json:"data" doc:"list of api endpoints to get data from"`
	Children []*Navigation `json:"children" doc:"Child navigation"`
}

func (self *Navigation) IsActive() bool {
	return self.Display.IsActive
}
func (self *Navigation) IsUri() bool {
	return self.Display.InUri
}
func (self *Navigation) IsHeader() bool {
	return self.Display.IsHeader
}

// New will create a standard Navigation struct using the name and uri passed
// It allows sthe more complex properties to be passed as variadic arguments
// and maakes use of a switch on type to assign them correctly.
// Allows for short hand creation and nest creation like:
//
//	New("test", "/test", New("foo", "/test/foo"))
func New(name string, uri string, others ...interface{}) (nav *Navigation) {
	nav = &Navigation{
		Name:     name,
		Uri:      uri,
		Display:  &Display{},
		Data:     []*Data{},
		Children: []*Navigation{},
	}
	// allow more complex parts of the nav item to be passed as variadic args
	for _, arg := range others {
		switch arg.(type) {
		case *Display:
			nav.Display = arg.(*Display)
		case *Data:
			nav.Data = append(nav.Data, arg.(*Data))
		case []*Data:
			nav.Data = arg.([]*Data)
		case *Navigation:
			nav.Children = append(nav.Children, arg.(*Navigation))
		case []*Navigation:
			nav.Children = arg.([]*Navigation)

		}
	}

	return
}
