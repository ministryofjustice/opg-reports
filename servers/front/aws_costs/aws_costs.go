package aws_costs

import (
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/shared/datarow"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/page"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/template"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/mw"

	"github.com/ministryofjustice/opg-reports/shared/convert"
)

const ytdTemplate string = "aws-costs-index"
const monthlyTaxTemplate string = "aws-costs-monthly-tax-totals"
const monthlyTemplate string = "aws-costs-monthly"

func decorators(re *aws_costs.ApiResponse, server *front.FrontServer, navItem *nav.Nav, r *http.Request) {
	re.Organisation = server.Config.Organisation
	re.PageTitle = navItem.Name
	if len(server.Config.Navigation) > 0 {
		top, active := nav.Level(server.Config.Navigation, r)
		re.NavigationActive = active
		re.NavigationTop = top
		re.NavigationSide = active.Navigation
	}
}

func rows(re *aws_costs.ApiResponse) {
	mapped, _ := convert.Maps(re.Result)
	intervals := map[string][]string{"interval": re.DateRange}
	values := map[string]string{"interval": "total"}
	re.Rows = datarow.DataRows(mapped, re.Columns, intervals, values)
}

// Handler
// Implements front.FrontHandlerFunc
func Handler(server *front.FrontServer, navItem *nav.Nav, w http.ResponseWriter, r *http.Request) {

	var (
		pg           = page.New[*aws_costs.ApiResponse](server, navItem)
		pageTemplate = template.New(navItem.Template, server.Templates, w)
	)
	responses, err := httphandler.GetAll(server.ApiSchema, server.ApiAddr, navItem.DataSources)

	if err != nil {
		slog.Error("error getting responses")
	}
	pg.SetNavigation(server, r)

	for key, handler := range responses {
		api, err := convert.UnmarshalR[*aws_costs.ApiResponse](handler.Response)
		if err != nil {
			return
		}
		if key != "ytd" {
			rows(api)
		}
		pg.Results[key] = api
	}

	pageTemplate.Run(pg)

}

func Register(mux *http.ServeMux, frontServer *front.FrontServer) {
	navigation := frontServer.Config.Navigation
	handledTemplates := []string{ytdTemplate, monthlyTaxTemplate, monthlyTemplate}

	for _, templateName := range handledTemplates {
		navItems := nav.ForTemplate(templateName, navigation)
		for _, navItem := range navItems {
			handler := front.Wrap(frontServer, navItem, Handler)

			slog.Info("[front] registering", slog.String("endpoint", "aws_costs"), slog.String("uri", navItem.Uri), slog.String("handler", "Handler"))
			mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
		}
	}

}
