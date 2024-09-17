package aws_uptime

import (
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/api/aws_uptime"
	"github.com/ministryofjustice/opg-reports/servers/shared/datarow"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/page"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/template"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"

	"github.com/ministryofjustice/opg-reports/shared/convert"
)

const uptimeTemplate string = "aws-uptime"

func rows(re *aws_uptime.ApiResponse) {
	mapped, _ := convert.Maps(re.Result)
	intervals := map[string][]string{"interval": re.DateRange}
	values := map[string]string{"interval": "average"}
	re.Rows = datarow.DataRows(mapped, re.Columns, intervals, values)
}

// Handler
// Implements front.FrontHandlerFunc
func Handler(server *front.FrontServer, navItem *nav.Nav, w http.ResponseWriter, r *http.Request) {

	var (
		pg           = page.New[*aws_uptime.ApiResponse](server, navItem)
		pageTemplate = template.New(navItem.Template, server.Templates, w)
	)
	responses, err := httphandler.GetAll(server.ApiSchema, server.ApiAddr, navItem.DataSources)

	if err != nil {
		slog.Error("error getting responses")
	}
	pg.SetNavigation(server, r)

	for key, handler := range responses {
		api, err := convert.UnmarshalR[*aws_uptime.ApiResponse](handler.Response)
		if err != nil {
			return
		}
		rows(api)
		pg.Results[key] = api
	}

	pageTemplate.Run(pg)

}

// Register
func Register(mux *http.ServeMux, frontServer *front.FrontServer) {
	navigation := frontServer.Config.Navigation
	handledTemplates := []string{uptimeTemplate}

	for _, templateName := range handledTemplates {
		navItems := nav.ForTemplate(templateName, navigation)

		for _, navItem := range navItems {
			frontServer.Register(mux, navItem, Handler)

		}
	}

}
