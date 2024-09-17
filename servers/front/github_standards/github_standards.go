package github_standards

import (
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/page"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/template"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

const templateName string = "github-standards"

// ListHandler
// Implements front.FrontHandlerFunc
func ListHandler(server *front.FrontServer, navItem *nav.Nav, w http.ResponseWriter, r *http.Request) {
	var (
		pg           *page.Page[*github_standards.GHSResponse] = page.New[*github_standards.GHSResponse](server, navItem)
		pageTemplate                                           = template.New(navItem.Template, server.Templates, w)
	)
	responses, err := httphandler.GetAll(server.ApiSchema, server.ApiAddr, navItem.DataSources)

	if err != nil {
		slog.Error("error getting responses")
	}
	pg.SetNavigation(server, r)

	for key, handler := range responses {
		api, err := convert.UnmarshalR[*github_standards.GHSResponse](handler.Response)
		if err != nil {
			return
		}
		pg.Results[key] = api
	}

	pageTemplate.Run(pg)

}

func Register(mux *http.ServeMux, frontServer *front.FrontServer) {
	navigation := frontServer.Config.Navigation

	navItems := nav.ForTemplate(templateName, navigation)
	for _, navItem := range navItems {
		frontServer.Register(mux, navItem, ListHandler)

	}

}
