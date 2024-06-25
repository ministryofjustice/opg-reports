package server

import (
	"html/template"
	"log/slog"
	"net/http"
	"opg-reports/services/front/cnf"
)

type FrontWebServer struct {
	templateFiles []string
	Config        *cnf.Config
	Nav           *Navigation
}

func (s *FrontWebServer) Write(w http.ResponseWriter, status int, tmpl *template.Template, name string, data any) {
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

func New(conf *cnf.Config, templates []string) *FrontWebServer {
	nav := &Navigation{tree: conf.Sections}
	return &FrontWebServer{
		templateFiles: templates,
		Config:        conf,
		Nav:           nav,
	}
}
