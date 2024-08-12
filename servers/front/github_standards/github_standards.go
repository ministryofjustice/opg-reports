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
	"github.com/ministryofjustice/opg-reports/servers/front/template_functions"
	"github.com/ministryofjustice/opg-reports/servers/front/write"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
)

const templateName string = "github-standards"

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	nav := conf.Navigation
	navItem := navigation.ForTemplate(templateName, nav)

	var list = func(w http.ResponseWriter, r *http.Request) {
		data := getter.Api(conf, navItem, r)
		// fmt.Println(convert.Printify(data))
		for k, _ := range data {
			fmt.Println(k)
		}

		// sort out navigation
		top, active := navigation.Level(conf.Navigation, r)
		data["NavigationTop"] = top
		data["NavigationSide"] = active.Navigation

		// -- template rendering!
		status := http.StatusOK
		t, err := template.New(navItem.Template).Funcs(template_functions.Funcs()).ParseFiles(templates...)
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
