package github_standards

import (
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/template"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/mw"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

const templateName string = "github-standards"

func decorators(re *github_standards.GHSResponse, server *front.FrontServer, navItem *nav.Nav, r *http.Request) {
	re.Organisation = server.Config.Organisation
	re.PageTitle = navItem.Name
	if len(server.Config.Navigation) > 0 {
		top, active := nav.Level(server.Config.Navigation, r)
		re.NavigationActive = active
		re.NavigationTop = top
		re.NavigationSide = active.Navigation
	}
}

// ListHandler
// Implements front.FrontHandlerFunc
func ListHandler(server *front.FrontServer, navItem *nav.Nav, w http.ResponseWriter, r *http.Request) {
	var (
		data         interface{}
		mapData      = map[string]interface{}{}
		apiSchema    = server.ApiSchema
		apiAddr      = server.ApiAddr
		paths        = navItem.DataSources
		pageTemplate = template.New(navItem.Template, server.Templates, w)
	)
	responses, err := httphandler.GetAll(apiSchema, apiAddr, paths)
	count := len(responses)

	if err != nil {
		slog.Error("error getting responses")
	}
	for key, handler := range responses {
		gh, err := convert.UnmarshalR[*github_standards.GHSResponse](handler.Response)
		if err != nil {
			return
		}
		decorators(gh, server, navItem, r)
		if count > 1 {
			mapData[key] = gh
			data = mapData
		} else {
			data = gh
		}
	}

	pageTemplate.Run(data)

}

func Register(mux *http.ServeMux, frontServer *front.FrontServer) {
	navigation := frontServer.Config.Navigation

	navItems := nav.ForTemplate(templateName, navigation)
	for _, navItem := range navItems {
		handler := front.Wrap(frontServer, navItem, ListHandler)

		slog.Info("[front] registering", slog.String("endpoint", "github_standards"), slog.String("uri", navItem.Uri), slog.String("handler", "ListHandler"))
		mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
	}

}
