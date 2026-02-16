package page

import (
	"net/http"
	"strings"
)

type PageInfo struct {
	Name         string `json:"name"`
	GovUKVersion string `json:"govuk_version"` // --govuk_version
	Signature    string `json:"signature"`     // --signature
}

// PageContent is a base data structure that is used on all pages
//
// More complex pages should use their own struct, but can use this
// as part of that struct to ensure base needs are covered
type PageContent struct {
	Name         string `json:"-"`             // Name is the core name of the front end, used in page title
	GovUKVersion string `json:"govuk_version"` // GovUKVersion is the version number (minus v) that we're using in this front end
	Signature    string `json:"signature"`     // Signature is combination of the semver & git commit

	request *http.Request `json:"-"`
	paths   []string      `json:"-"`
}

// RequestPath returns the segment of the current request
func (self *PageContent) RequestPath(i int) (v string) {
	v = ""
	if len(self.paths) > i {
		v = self.paths[i]
	}
	return

}

// NewContent
func NewContent(request *http.Request, info *PageInfo) (pg PageContent) {
	var (
		uri   string   = strings.TrimPrefix(request.URL.String(), "/")
		paths []string = strings.Split(uri, "/")
	)
	return PageContent{
		Name:         info.Name,
		GovUKVersion: info.GovUKVersion,
		Signature:    info.Signature,
		request:      request,
		paths:        paths,
	}
}
