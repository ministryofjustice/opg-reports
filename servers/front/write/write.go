package write

import (
	"html/template"
	"log/slog"
	"net/http"
)

func Out(w http.ResponseWriter, status int, tmpl *template.Template, name string, data any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/html")
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		slog.Error("front template execute failed",
			slog.String("templateName", name),
			slog.String("error", err.Error()),
		)
	}
}
