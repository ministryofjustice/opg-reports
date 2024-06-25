package server

import (
	"html/template"
	"log/slog"
	"net/http"
)

func (s *FrontWebServer) Static(w http.ResponseWriter, r *http.Request) {
	active := s.Nav.Active(r)
	data := s.Nav.Data(r)
	templateFiles := s.templateFiles

	t, err := template.New(active.TemplateName).ParseFiles(templateFiles...)
	if err != nil {
		slog.Error("static handler template failure",
			slog.String("error", err.Error()),
			slog.String("uri", r.RequestURI),
			slog.String("active.Href", active.Href),
			slog.String("active.Template", active.TemplateName))
	}

	data["Organisation"] = s.Config.Organisation
	data["PageTitle"] = active.Name + " - "

	slog.Info("static handler",
		slog.String("uri", r.RequestURI),
		slog.String("active.Href", active.Href),
		slog.String("active.Template", active.TemplateName))

	s.Write(w, 200, t, active.TemplateName, data)
}
