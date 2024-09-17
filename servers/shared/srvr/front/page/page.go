package page

import (
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/response"
)

type Page[T response.Respond] struct {
	Organisation string
	PageTitle    string

	NavigationTop    map[string]*nav.Nav
	NavigationSide   []*nav.Nav
	NavigationActive *nav.Nav

	Rows    map[string]map[string]interface{}
	Results map[string]T
}

func (pg *Page[T]) SetNavigation(server *front.FrontServer, request *http.Request) {
	if len(server.Config.Navigation) > 0 {
		top, active := nav.Level(server.Config.Navigation, request)
		pg.NavigationActive = active
		pg.NavigationTop = top
		pg.NavigationSide = active.Navigation
	}
}

func (pg *Page[T]) Get(key string) T {
	return pg.Results[key]
}

func New[T response.Respond](server *front.FrontServer, navItem *nav.Nav) *Page[T] {
	return &Page[T]{
		Organisation: server.Config.Organisation,
		PageTitle:    navItem.Name,
		Results:      map[string]T{},
	}
}
