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

func decorators(re *aws_costs.CostResponse, conf *config.Config, navItem *navigation.NavigationItem, r *http.Request) {
	re.Organisation = conf.Organisation
	re.PageTitle = navItem.Name
	if len(conf.Navigation) > 0 {
		top, active := navigation.Level(conf.Navigation, r)
		re.NavigationActive = active
		re.NavigationTop = top
		re.NavigationSide = active.Navigation
	}
}

func rows(re *aws_costs.CostResponse) {
	mapped, _ := convert.Maps(re.Result)
	intervals := map[string][]string{"interval": re.DateRange}
	values := map[string]string{"interval": "total"}
	re.Rows = datarow.DataRows(mapped, re.Columns, intervals, values)
}

// YtdHandler
func YtdHandler(w http.ResponseWriter, r *http.Request, templates []string, conf *config.Config, navItem *navigation.NavigationItem) {
	var data interface{}
	mapData := map[string]interface{}{}
	if responses, err := getter.ApiHttpResponses(navItem, r); err == nil {
		count := len(responses)
		for key, rep := range responses {
			ytd, err := convert.UnmarshalR[*aws_costs.CostResponse](rep)
			if err != nil {
				return
			}
			// set the nav and org details
			decorators(ytd, conf, navItem, r)
			if count > 1 {
				mapData[key] = ytd
				data = mapData
			} else {
				data = ytd
			}
		}
	}
	helpers.OutputHandler(templates, navItem.Template, data, w)
}

// StandardHandler
func StandardHandler(w http.ResponseWriter, r *http.Request, templates []string, conf *config.Config, navItem *navigation.NavigationItem) {
	var data interface{}
	mapData := map[string]interface{}{}

	if responses, err := getter.ApiHttpResponses(navItem, r); err == nil {
		count := len(responses)
		for key, rep := range responses {
			mtx, err := convert.UnmarshalR[*aws_costs.CostResponse](rep)
			if err != nil {
				return
			}
			// -- create rows from the response
			decorators(mtx, conf, navItem, r)
			rows(mtx)

			if count > 1 {
				mapData[key] = mtx
				data = mapData
			} else {
				data = mtx
			}
		}
	}

	helpers.OutputHandler(templates, navItem.Template, data, w)

}

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	nav := conf.Navigation

	// -- year to date
	ytds := navigation.ForTemplateList(ytdTemplate, nav)
	for _, navItem := range ytds {
		var handler = func(w http.ResponseWriter, r *http.Request) {
			YtdHandler(w, r, templates, conf, navItem)
		}
		slog.Info("[front] register", slog.String("endpoint", "aws_costs"), slog.String("ytd", navItem.Uri))
		mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
	}

	// -- monthly totals with tax split
	taxes := navigation.ForTemplateList(monthlyTaxTemplate, nav)
	for _, navItem := range taxes {
		var handler = func(w http.ResponseWriter, r *http.Request) {
			StandardHandler(w, r, templates, conf, navItem)
		}
		slog.Info("[front] register", slog.String("endpoint", "aws_costs"), slog.String("monthly tax", navItem.Uri))
		mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
	}

	// -- list by month
	months := navigation.ForTemplateList(monthlyTemplate, nav)
	for _, navItem := range months {
		var handler = func(w http.ResponseWriter, r *http.Request) {
			StandardHandler(w, r, templates, conf, navItem)
		}
		slog.Info("[front] register", slog.String("endpoint", "aws_costs"), slog.String("monthly", navItem.Uri))
		mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
	}

}
