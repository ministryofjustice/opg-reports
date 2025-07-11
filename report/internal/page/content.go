package page

import (
	"net/http"
	"opg-reports/report/config"
	"strings"
)

// PageContent is a base data structure that is used on all pages
//
// More complex pages should use their own struct, but can use this
// as part of that struct to ensure base needs are covered
type PageContent struct {
	Name         string // Name is the core name of the front end, used in page title
	GovUKVersion string // GovUKVersion is the version number (minus v) that we're using in this front end
	Signature    string // Signature is combination of the semver & git commit
	RequestPath0 string // First segment of Request.URL.String()
	RequestPath1 string // Second segment of Request.URL.String()

	Teams []string
}

func DefaultContent(conf *config.Config, request *http.Request) (pg PageContent) {
	paths := strings.Split(request.URL.String(), "/")
	pg = PageContent{
		Name:         conf.Servers.Front.Name,
		GovUKVersion: strings.TrimPrefix(conf.GovUK.Front.ReleaseTag, "v"),
		Signature:    conf.Versions.Signature(),
		RequestPath0: paths[0],
	}
	if len(paths) > 0 {
		pg.RequestPath1 = paths[1]
	}
	return
}
