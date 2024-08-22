package github_standards

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/front/template_helpers"
	"github.com/ministryofjustice/opg-reports/servers/front/write"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
)

const templateName string = "github-standards"

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	nav := conf.Navigation
	navItem := navigation.ForTemplate(templateName, nav)

	var list = func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{"Result": nil}

		// if theres no error, then process the response normally
		if responses, err := getter.ApiResponses(conf, navItem, r); err == nil {
			data = getter.ParseApiResponse(responses)
			metadata := data["Metadata"].(map[string]interface{})
			counters := metadata["counters"].(map[string]interface{})
			this := counters["this"].(map[string]interface{})
			total := (this["count"].(float64))
			base := (this["compliant_baseline"].(float64))
			ext := (this["compliant_extended"].(float64))
			percent := base / (total / 100)

			data["Total"] = total
			data["PassedBaseline"] = base
			data["PassedExtended"] = ext
			data["Percentage"] = fmt.Sprintf("%.2f", percent)

		} else {
			slog.Error("had error back from api getter", slog.String("err", err.Error()))
		}
		data["Organisation"] = conf.Organisation
		data["PageTitle"] = navItem.Name + " - "
		// sort out navigation
		top, active := navigation.Level(conf.Navigation, r)
		data["NavigationTop"] = top
		data["NavigationSide"] = active.Navigation
		// -- template rendering!
		status := http.StatusOK
		t, err := template.New(navItem.Template).Funcs(template_helpers.Funcs()).ParseFiles(templates...)
		if err != nil {
			slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
			status = http.StatusBadGateway
		}
		write.Out(w, status, t, templateName, data)
	}
	// -- register
	slog.Info("[front] register",
		slog.String("handler", "github_standards_list"),
		slog.String("uri", navItem.Uri))

	mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(list, mw.Logging, mw.SecurityHeaders))
}
