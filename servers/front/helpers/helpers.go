package helpers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/front/template_helpers"
	"github.com/ministryofjustice/opg-reports/servers/front/write"
)

func DataCleanup(data map[string]interface{}, conf *config.Config, navItem *navigation.NavigationItem, r *http.Request) map[string]interface{} {
	data["Organisation"] = conf.Organisation
	data["PageTitle"] = navItem.Name
	// sort out navigation
	if len(conf.Navigation) > 0 {
		top, active := navigation.Level(conf.Navigation, r)
		data["NavigationTop"] = top
		data["NavigationSide"] = active.Navigation
	} else {
		data["NavigationTop"] = map[string]*navigation.NavigationItem{}
		data["NavigationSide"] = []*navigation.NavigationItem{}
	}
	return data
}

func OutputHandler(templates []string, templateName string, data map[string]interface{}, w http.ResponseWriter) {
	status := http.StatusOK
	t, err := template.New(templateName).Funcs(template_helpers.Funcs()).ParseFiles(templates...)
	if err != nil {
		slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
		status = http.StatusBadGateway
	}
	write.Out(w, status, t, templateName, data)
}
