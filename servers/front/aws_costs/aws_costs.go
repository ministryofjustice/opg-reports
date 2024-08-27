package aws_costs

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/front/helpers"
	"github.com/ministryofjustice/opg-reports/servers/shared/datarow"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

const ytdTemplate string = "aws-costs-index"
const monthlyTaxTemplate string = "aws-costs-monthly-tax-totals"
const monthlyTemplate string = "aws-costs-monthly"

func decorators(re *aws_costs.ApiResponse, conf *config.Config, navItem *navigation.NavigationItem, r *http.Request) {
	re.Organisation = conf.Organisation
	re.PageTitle = navItem.Name
	if len(conf.Navigation) > 0 {
		top, active := navigation.Level(conf.Navigation, r)
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

func Handler(w http.ResponseWriter, r *http.Request, templates []string, conf *config.Config, navItem *navigation.NavigationItem, boundTemplate string) {
	var data interface{}
	mapData := map[string]interface{}{}
	if responses, err := getter.ApiHttpResponses(navItem, r); err == nil {
		count := len(responses)
		for key, rep := range responses {
			api, err := convert.UnmarshalR[*aws_costs.ApiResponse](rep)
			if err != nil {
				return
			}
			// set the nav and org details
			decorators(api, conf, navItem, r)
			if boundTemplate != ytdTemplate {
				rows(api)
			}
			if count > 1 {
				mapData[key] = api
				data = mapData
			} else {
				data = api
			}
		}
	}
	helpers.OutputHandler(templates, navItem.Template, data, w)
}

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	var (
		nav              = conf.Navigation
		handledTemplates = []string{ytdTemplate, monthlyTaxTemplate, monthlyTemplate}
	)

	for _, template := range handledTemplates {
		templateNavs := navigation.ForTemplateList(template, nav)
		for _, navItem := range templateNavs {
			var handler = func(w http.ResponseWriter, r *http.Request) {
				Handler(w, r, templates, conf, navItem, template)
			}
			slog.Info("[front] register", slog.String("endpoint", "aws_costs"), slog.String("uri", navItem.Uri), slog.String("template", template))
			mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
		}

	}

}
