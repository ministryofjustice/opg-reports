package front

import "opg-reports/report/packages/httpx"

// Page is the html home page; extends the response data
type Page struct {
	httpx.HTMLResponseData
}

// TemplateName returns the html template name to use
func (self *Page) TemplateName() string {
	return "landing-page"
}

// SetTeams
func (self *Page) SetTeams(teams []string) {
	self.Teams = teams
}
