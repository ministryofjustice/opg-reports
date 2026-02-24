package htmlpage

import (
	"net/http"
	"strings"
)

type Args struct {
	Title        string // page title (<title>)
	Name         string // page name - used in heading
	GovUKVersion string // the govuk version (no v)
	SemVer       string // sematic version

}

type HTMLPage struct {
	Title        string   // page title (<title>)
	Name         string   // Name is used for page title
	GovUKVersion string   // GovUKVersion is the version number (minus v) that we're using in this front end
	SemVer       string   // sematic version
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
		Title:        args.Title,
		Name:         args.Name,
		GovUKVersion: args.GovUKVersion,
		SemVer:       args.SemVer,
		request:      request,
		paths:        paths,
	}
}
