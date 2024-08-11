package github_standards

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
)

const template string = "github-standards"

func Register(ctx context.Context, mux *http.ServeMux, conf *config.Config, templates []string) {
	nav := conf.Navigation
	navItem := navigation.ForTemplate(template, nav)

	var list = func(w http.ResponseWriter, r *http.Request) {
		getter.Api(conf, navItem, r)
		// -- template rendering!
		// t, err := template.New(navItem.Template).Funcs(tmpl.Funcs()).ParseFiles(templates...)
		// if err != nil {
		// 	slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
		// 	return
		// }
	}
	// -- register
	slog.Info("[front] register",
		slog.String("handler", "github_standards_list"),
		slog.String("uri", navItem.Uri))
	mux.HandleFunc(navItem.Uri+"{$}", mw.Middleware(list, mw.Logging, mw.SecurityHeaders))
}
