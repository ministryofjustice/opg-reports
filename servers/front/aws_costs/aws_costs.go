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
)

const monthlyTaxTemplateName string = "aws-costs-monthly-tax"

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	nav := conf.Navigation
	monthlyTaxNavItem := navigation.ForTemplate(monthlyTaxTemplateName, nav)

	var monthlyTax = func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{"Result": nil, "Months": []string{}}

		// if theres no error, then process the response normally
		if apiData, err := getter.Api(conf, monthlyTaxNavItem, r); err == nil {
			data = apiData
			metadata := data["Metadata"].(map[string]interface{})
			data["Months"] = metadata["months"].([]interface{})
			data["Columns"] = metadata["columns"].(map[string]interface{})

		} else {
			slog.Error("had error back from api getter", slog.String("err", err.Error()))
		}
		data["Organisation"] = conf.Organisation
		data["PageTitle"] = monthlyTaxNavItem.Name + " "
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
		slog.String("monthlyTaxNavItem", monthlyTaxNavItem.Uri))

	mux.HandleFunc(monthlyTaxNavItem.Uri+"{$}", mw.Middleware(monthlyTax, mw.Logging, mw.SecurityHeaders))
}
