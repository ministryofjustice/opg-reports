// Package navigation used for navigation config on the front end.
//
// Allows for mapping a front end url to a series of data sources
// on the api and what template to render for each page on the
// front end.
package navigation

import (
	"slices"
	"strings"

	"github.com/ministryofjustice/opg-reports/internal/endpoints"
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

// ResponseTransformer is used to handle processing and changing the
// api data before display
// The function should process the body and replace any values
// and return the updated body as the result
type ResponseTransformer func(body interface{}) interface{}

type Data struct {
	Source      endpoints.ApiEndpoint `json:"source" faker:"uri" doc:"API data source location"`
	Namespace   string                `json:"namepsace" faker:"word" doc:"Namespace the Result should be copied over into for front end parsing"`
	Body        interface{}           `json:"body" doc:"A pointer to the struct that this data would use as the body content."`
	Transformer ResponseTransformer   `json:"-" doc:"Function to transform the api data into front end structure"`
}

func NewData(source endpoints.ApiEndpoint, namespace string) *Data {
	return &Data{Source: source, Namespace: namespace}
}

// Navigation represents a tree structured navigation that is rendered on the
// front end of website to show sections at a time
// It also contains details on how to fetch data for each of the pages it represents
type Navigation struct {
	Name    string   `json:"name" doc:"Name is used when displaying this navigation link."`
	Uri     string   `json:"uri" doc:"Uri is front end url this navigation item will render for."`
	Display *Display `json:"display" doc:"rendering related details"`
	Data    []*Data  `json:"data" doc:"list of api endpoints to get data from"`

	children []*Navigation
	parent   *Navigation
}

// IsActive returns the Display.IsActive value
// - this makes it easier within templates to check
func (self *Navigation) IsActive() bool {
	return self.Display.IsActive
}

// InUri returns the Display.InUri value
// - this makes it easier within templates to check
func (self *Navigation) InUri() bool {
	return self.Display.InUri
}

func (self *Navigation) Matches(other *Navigation) bool {
	if other == nil {
		return false
	}
	return self.Uri == other.Uri
}

// IsHeader returns the Display.IsHeader value
// - this makes it easier within templates to check
func (self *Navigation) IsHeader() bool {
	return self.Display.IsHeader
}

// IsActiveOrInUri is used in the top bar nav
// to check it the item should be marked as active
func (self *Navigation) IsActiveOrInUri() bool {
	if self.Uri == "/" {
		return self.IsActive()
	}
	return self.IsActive() || self.InUri()
}

// Parent returns the .parent navigation item
// for this - or nil
func (self *Navigation) Parent() *Navigation {
	return self.parent
}

// Children returns a slice of navigation items that
// are the child of this
func (self *Navigation) Children() []*Navigation {
	return self.children
}

// AddChild add this kids passed as being the children
// of this node - updates their internal .parent
// as well
func (self *Navigation) AddChild(kids ...*Navigation) {
	for _, child := range kids {
		child.parent = self
		self.children = append(self.children, child)
	}
}

func (self *Navigation) Names() (names string) {
	var sep string = " - "
	path := RootPath(self)
	if len(path) > 0 {
		list := []string{}
		for _, p := range path {
			list = append(list, p.Name)
		}
		slices.Reverse(list)
		names = strings.Join(list, sep)
	}

	return
}

// Root finds the top level node from the node passed,
// recursing up the tree by using hte .parent
func Root(node *Navigation) (p *Navigation) {
	p = node
	// more parents, so go up another level
	if parent := node.Parent(); parent != nil {
		p = Root(parent)
	}
	return

}

func RootPath(node *Navigation) (nodes []*Navigation) {
	nodes = []*Navigation{node}
	if parent := node.Parent(); parent != nil {
		nodes = append(nodes, RootPath(parent)...)
	}
	return

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
		children: []*Navigation{},
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
			nav.AddChild(arg.(*Navigation))
		case []*Navigation:
			kids := arg.([]*Navigation)
			nav.AddChild(kids...)

		}
	}

	return
}
