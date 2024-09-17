package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
)

type Page struct {
	Organisation string
	PageTitle    string

	NavigationTop    map[string]*nav.Nav
	NavigationSide   []*nav.Nav
	NavigationActive *nav.Nav

	Rows    map[string]map[string]interface{}
	Results map[string]interface{}
}

func (pg *Page) SetNavigation(server *FrontServer, request *http.Request) {
	if len(server.Config.Navigation) > 0 {
		top, active := nav.Level(server.Config.Navigation, request)
		pg.NavigationActive = active
		pg.NavigationTop = top
		pg.NavigationSide = active.Navigation
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
