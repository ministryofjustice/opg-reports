package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
)

type Page struct {
	Organisation string
	PageTitle    string

	NavigationTopbarItems  map[string]*nav.Nav
	ActiveSection          *nav.Nav
	NavigationSidebarItems []*nav.Nav
	CurrentPage            *nav.Nav

	Results map[string]interface{}
}

func (pg *Page) SetNavigation(server *FrontServer, request *http.Request) {
	if len(server.Config.Navigation) > 0 {
		top, activeSection, activePage := nav.Level(server.Config.Navigation, request)
		pg.ActiveSection = activeSection
		pg.CurrentPage = activePage
		pg.NavigationTopbarItems = top

		if activeSection != nil {
			pg.NavigationSidebarItems = activeSection.Navigation
		}

	}
}

func (pg *Page) Get(key string) interface{} {
	return pg.Results[key]
}

func NewPage(server *FrontServer, navItem *nav.Nav, request *http.Request) (pg *Page) {
	pg = &Page{
		Organisation: server.Config.Organisation,
		PageTitle:    navItem.Name,
		Results:      server.PageData,
	}
	if request != nil {
		pg.SetNavigation(server, request)
	}
	return
}
