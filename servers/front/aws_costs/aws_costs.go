package aws_costs

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"slices"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/front/rows"
	"github.com/ministryofjustice/opg-reports/servers/front/template_helpers"
	"github.com/ministryofjustice/opg-reports/servers/front/write"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

const ytdTemplate string = "aws-costs-index"
const monthlyTaxTemplate string = "aws-costs-monthly-tax-totals"
const monthlyTemplate string = "aws-costs-monthly"

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	nav := conf.Navigation

	// -- year to date
	ytds := navigation.ForTemplateList(ytdTemplate, nav)
	for _, navItem := range ytds {
		var handler = func(w http.ResponseWriter, r *http.Request) {
			data := map[string]interface{}{"Result": nil}
			if apiData, err := getter.Api(conf, navItem, r); err == nil {
				data = apiData
				// total
				res := data["Result"].([]interface{})
				first := res[0].(map[string]interface{})
				data["Total"] = first["total"].(float64)
			}
			data = dataCleanup(data, conf, navItem, r)
			outputHandler(templates, navItem.Template, data, w)
		}

		slog.Info("[front] register", slog.String("endpoint", "aws_costs"), slog.String("ytd", navItem.Uri))
		mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
	}

	// -- monthly totals with tax split
	taxes := navigation.ForTemplateList(monthlyTaxTemplate, nav)
	for _, navItem := range taxes {
		var handler = func(w http.ResponseWriter, r *http.Request) {
			data := map[string]interface{}{"Result": nil}
			if apiData, err := getter.Api(conf, navItem, r); err == nil {
				data = apiData
			}
			data = dataCleanup(data, conf, navItem, r)
			outputHandler(templates, navItem.Template, data, w)
		}

		slog.Info("[front] register", slog.String("endpoint", "aws_costs"), slog.String("monthly tax", navItem.Uri))
		mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
	}

	// -- list by month
	months := navigation.ForTemplateList(monthlyTemplate, nav)
	for _, navItem := range months {
		var handler = func(w http.ResponseWriter, r *http.Request) {
			data := map[string]interface{}{"Result": nil}
			if apiData, err := getter.Api(conf, navItem, r); err == nil {
				data = apiData
				// -- get date ranges to use for mapping
				dataRange := data["DateRange"].([]string)
				// -- get detailed columns
				columns := data["ColumnsDetailed"].(map[string][]interface{})
				// -- map to rows from the data set
				intervals := map[string][]string{"interval": dataRange}
				values := map[string]string{"interval": "total"}
				res := data["Result"].([]interface{})
				data["Result"] = rows.DataRows(res, columns, intervals, values)
				// -- need to create the filters for this version
				filters := []string{}
				for i, col := range data["ColumnsOrdered"].([]string) {
					filters = append(filters, fmt.Sprintf("%d.%s", i+1, convert.Title(col)))
				}
				slices.Sort(filters)
				data["Filters"] = filters
			}
			data = dataCleanup(data, conf, navItem, r)
			outputHandler(templates, navItem.Template, data, w)
		}

		slog.Info("[front] register", slog.String("endpoint", "aws_costs"), slog.String("monthly", navItem.Uri))
		mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
	}

}

func dataCleanup(data map[string]interface{}, conf *config.Config, navItem *navigation.NavigationItem, r *http.Request) map[string]interface{} {
	data["Organisation"] = conf.Organisation
	data["PageTitle"] = navItem.Name + " - "
	// sort out navigation
	top, active := navigation.Level(conf.Navigation, r)
	data["NavigationTop"] = top
	data["NavigationSide"] = active.Navigation
	return data
}

func outputHandler(templates []string, templateName string, data map[string]interface{}, w http.ResponseWriter) {
	status := http.StatusOK
	t, err := template.New(ytdTemplate).Funcs(template_helpers.Funcs()).ParseFiles(templates...)
	if err != nil {
		slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
		status = http.StatusBadGateway
	}
	write.Out(w, status, t, templateName, data)
}
