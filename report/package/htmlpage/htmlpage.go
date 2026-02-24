package htmlpage

import (
	"net/http"
	"strings"
)

type Args struct {
	Name         string
	GovUKVersion string
}

type HTMLPage struct {
	Name         string   // Name is used for page title
	GovUKVersion string   // GovUKVersion is the version number (minus v) that we're using in this front end
	Teams        []string // Teams are used for the page navigation

	request *http.Request
	paths   []string
}

// RequestPath returns the segment of the current request
func (self *HTMLPage) RequestPath(i int) (v string) {
	v = ""
	if len(self.paths) > i {
		v = self.paths[i]
	}
	return

}

// New
func New(request *http.Request, args *Args) (pg HTMLPage) {
	var (
		uri   string   = strings.TrimPrefix(request.URL.Path, "/")
		paths []string = strings.Split(uri, "/")
	)
	return HTMLPage{
		Name:         args.Name,
		GovUKVersion: args.GovUKVersion,
		request:      request,
		paths:        paths,
	}
}
