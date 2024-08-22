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

func ListHandler(w http.ResponseWriter, r *http.Request, templates []string, conf *config.Config, navItem *navigation.NavigationItem) {
	data := map[string]interface{}{"Result": nil}
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

	}
	data = dataCleanup(data, conf, navItem, r)
	outputHandler(templates, navItem.Template, data, w)
}

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	nav := conf.Navigation
	navItems := navigation.ForTemplateList(templateName, nav)
	for _, navItem := range navItems {
		var handler = func(w http.ResponseWriter, r *http.Request) {
			ListHandler(w, r, templates, conf, navItem)
		}
		slog.Info("[front] register", slog.String("endpoint", "githug_standards"), slog.String("list", navItem.Uri))
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
	t, err := template.New(templateName).Funcs(template_helpers.Funcs()).ParseFiles(templates...)
	if err != nil {
		slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
		status = http.StatusBadGateway
	}
	write.Out(w, status, t, templateName, data)
}
