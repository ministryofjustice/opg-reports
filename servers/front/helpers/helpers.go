package helpers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/front/template_helpers"
	"github.com/ministryofjustice/opg-reports/servers/front/write"
)

func OutputHandler(templates []string, templateName string, data any, w http.ResponseWriter) {
	status := http.StatusOK
	t, err := template.New(templateName).Funcs(template_helpers.Funcs()).ParseFiles(templates...)
	if err != nil {
		slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
		status = http.StatusBadGateway
	}
	write.Out(w, status, t, templateName, data)
}
