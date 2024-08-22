package aws_costs

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/front/front_templates"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/front/rows"
	"github.com/ministryofjustice/opg-reports/servers/front/write"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

const ytdTemplate string = "aws-costs-index"
const monthlyTaxTemplate string = "aws-costs-monthly-tax-totals"
const monthlyTemplate string = "aws-costs-monthly"

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
	t, err := template.New(ytdTemplate).Funcs(front_templates.Funcs()).ParseFiles(templates...)
	if err != nil {
		slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
		status = http.StatusBadGateway
	}
	write.Out(w, status, t, templateName, data)
}

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	nav := conf.Navigation

	// -- year to date
	ytds := navigation.ForTemplateList(ytdTemplate, nav)
	for _, navItem := range ytds {
		var handler = func(w http.ResponseWriter, r *http.Request) {
			data := map[string]interface{}{"Result": nil}
			if apiData, err := getter.Api(conf, navItem, r); err == nil {
				data = apiData
				metadata := data["Metadata"].(map[string]interface{})
				// start & end date
				sd := metadata["start_date"].(interface{})
				ed := metadata["end_date"].(interface{})
				data["StartDate"] = dates.Time(sd.(string))
				data["EndDate"] = dates.Time(ed.(string))
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
				metadata := data["Metadata"].(map[string]interface{})
				data["DateRange"] = metadata["date_range"].([]interface{})
				data["Columns"] = metadata["columns"].(map[string]interface{})
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
				metadata := data["Metadata"].(map[string]interface{})
				data["DateRange"] = metadata["date_range"].([]interface{})

				// -- convert data ranges to strings
				dataRange := []string{}
				for _, dr := range metadata["date_range"].([]interface{}) {
					dataRange = append(dataRange, dr.(string))
				}
				// -- convert columns
				cols := metadata["columns"].(map[string]interface{})
				columns := map[string][]interface{}{}
				for col, val := range cols {
					columns[col] = val.([]interface{})
				}
				// column ordering
				colNames := []string{}
				for _, col := range metadata["column_ordering"].([]interface{}) {
					colNames = append(colNames, col.(string))
				}
				data["Columns"] = colNames

				intervals := map[string][]string{"interval": dataRange}
				values := map[string]string{"interval": "total"}
				res := data["Result"].([]interface{})
				data["Result"] = rows.DataToRows(res, columns, intervals, values)

			}
			data = dataCleanup(data, conf, navItem, r)
			outputHandler(templates, navItem.Template, data, w)
		}

		slog.Info("[front] register", slog.String("endpoint", "aws_costs"), slog.String("monthly", navItem.Uri))
		mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(handler, mw.Logging, mw.SecurityHeaders))
	}

}
