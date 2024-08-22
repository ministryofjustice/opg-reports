package aws_costs

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/front/helpers"
	"github.com/ministryofjustice/opg-reports/servers/front/rows"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

const ytdTemplate string = "aws-costs-index"
const monthlyTaxTemplate string = "aws-costs-monthly-tax-totals"
const monthlyTemplate string = "aws-costs-monthly"

// YtdHandler
func YtdHandler(w http.ResponseWriter, r *http.Request, templates []string, conf *config.Config, navItem *navigation.NavigationItem) {
	data := map[string]interface{}{"Result": nil}
	if responses, err := getter.ApiResponses(navItem, r); err == nil {
		data = getter.ParseApiResponse(responses)
		// total
		res := data["Result"].([]interface{})
		first := res[0].(map[string]interface{})
		data["Total"] = first["total"].(float64)
	}
	data = helpers.DataCleanup(data, conf, navItem, r)
	helpers.OutputHandler(templates, navItem.Template, data, w)
}

// MonthlyTaxHandler
func MonthlyTaxHandler(w http.ResponseWriter, r *http.Request, templates []string, conf *config.Config, navItem *navigation.NavigationItem) {
	data := map[string]interface{}{"Result": nil}
	if responses, err := getter.ApiResponses(navItem, r); err == nil {
		data = getter.ParseApiResponse(responses)
	}
	data = helpers.DataCleanup(data, conf, navItem, r)
	helpers.OutputHandler(templates, navItem.Template, data, w)

}

// StandardHandler
func StandardHandler(w http.ResponseWriter, r *http.Request, templates []string, conf *config.Config, navItem *navigation.NavigationItem) {
	data := map[string]interface{}{"Result": nil}
	if responses, err := getter.ApiResponses(navItem, r); err == nil {
		data = getter.ParseApiResponse(responses)
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
	data = helpers.DataCleanup(data, conf, navItem, r)
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
			MonthlyTaxHandler(w, r, templates, conf, navItem)
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
