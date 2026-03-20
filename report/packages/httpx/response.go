package httpx

import "strings"

// ResponseData is a default struct that exposes
// version & request data which can then be
// expanded on by other structs.
type ResponseData struct {
	Version      string       `json:"version"`
	Request      *RequestData `json:"request"`
	GovUKVersion string       `json:"govuk_version"`
	// Teams is list of team names used within the site navigation
	Teams []string `json:"-"`
}

// TemplateName returns empty by default
func (self *ResponseData) TemplateName() string {
	return ""
}

// SetVersion pushes the version string into the response object
func (self *ResponseData) SetVersion(v string) {
	self.Version = v
}

// SetGovukVersion pushes the gov uk version into the setup
func (self *ResponseData) SetGovukVersion(v string) {
	self.GovUKVersion = v
}

// SetRequestData pushes the processed incoming request data
// back on to the response
func (self *ResponseData) SetRequestData(rd *RequestData) {
	self.Request = rd
}

// SetTeams
func (self *ResponseData) SetTeams(teams []string) {
	self.Teams = teams
}

// HTMLResponseData is extension of ResponseData that has some
// additional functionality for the html pages that are
// generated
type HTMLResponseData struct {
	ResponseData
	// used page title (<title>)
	Title string
	// Name is used for heading
	Name string
	// used to store request url in chunks for easier use in front end navigation
	paths []string
}

// RequestPath returns the segment of the current request
func (self *HTMLResponseData) RequestPath(i int) (v string) {
	var uri string = strings.Trim(self.Request.Request().URL.Path, "/")
	v = ""
	// set paths if its empty
	if len(self.paths) == 0 {
		self.paths = strings.Split(uri, "/")
	}
	// find the path value
	if len(self.paths) > i {
		v = self.paths[i]
	}
	return

}
