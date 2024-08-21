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
	"github.com/ministryofjustice/opg-reports/servers/front/write"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

const ytdTemplate string = "aws-costs-index"
const monthlyTaxTemplateName string = "aws-costs-monthly-totals"

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	nav := conf.Navigation

	// -- ytd page
	ytdNav := navigation.ForTemplate(ytdTemplate, nav)
	var ytd = func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{"Result": nil}

		// if theres no error, then process the response normally
		if apiData, err := getter.Api(conf, ytdNav, r); err == nil {
			data = apiData
			metadata := data["Metadata"].(map[string]interface{})
			sd := metadata["start_date"].(interface{})
			ed := metadata["end_date"].(interface{})
			data["StartDate"] = dates.Time(sd.(string))
			data["EndDate"] = dates.Time(ed.(string))

			// total
			res := data["Result"].([]interface{})
			first := res[0].(map[string]interface{})
			data["Total"] = first["total"].(float64)
		} else {
			slog.Error("had error back from api getter", slog.String("err", err.Error()))
		}
		data["Organisation"] = conf.Organisation
		data["PageTitle"] = ytdNav.Name + " - "
		// sort out navigation
		top, active := navigation.Level(conf.Navigation, r)
		data["NavigationTop"] = top
		data["NavigationSide"] = active.Navigation
		// -- template rendering!
		status := http.StatusOK
		t, err := template.New(ytdNav.Template).Funcs(front_templates.Funcs()).ParseFiles(templates...)
		if err != nil {
			slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
			status = http.StatusBadGateway
		}
		write.Out(w, status, t, ytdTemplate, data)
	}

	// -- monthly tax breakdown
	monthlyTaxNavItem := navigation.ForTemplate(monthlyTaxTemplateName, nav)
	var monthlyTax = func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{"Result": nil, "DateRange": []string{}}

		// if theres no error, then process the response normally
		if apiData, err := getter.Api(conf, monthlyTaxNavItem, r); err == nil {
			data = apiData
			metadata := data["Metadata"].(map[string]interface{})
			data["DateRange"] = metadata["date_range"].([]interface{})
			data["Columns"] = metadata["columns"].(map[string]interface{})

		} else {
			slog.Error("had error back from api getter", slog.String("err", err.Error()))
		}
		data["Organisation"] = conf.Organisation
		data["PageTitle"] = monthlyTaxNavItem.Name + " - "
		// sort out navigation
		top, active := navigation.Level(conf.Navigation, r)
		data["NavigationTop"] = top
		data["NavigationSide"] = active.Navigation
		// -- template rendering!
		status := http.StatusOK
		t, err := template.New(monthlyTaxNavItem.Template).Funcs(front_templates.Funcs()).ParseFiles(templates...)
		if err != nil {
			slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
			status = http.StatusBadGateway
		}
		write.Out(w, status, t, monthlyTaxTemplateName, data)
	}

	// -- register
	slog.Info("[front] register",
		slog.String("handler", "aws_costs"),
		slog.String("monthlyTaxNavItem", monthlyTaxNavItem.Uri),
		slog.String("ytd", ytdNav.Uri))

	mux.HandleFunc(monthlyTaxNavItem.Uri+"{$}", mw.Middleware(monthlyTax, mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(ytdNav.Uri+"{$}", mw.Middleware(ytd, mw.Logging, mw.SecurityHeaders))
}
