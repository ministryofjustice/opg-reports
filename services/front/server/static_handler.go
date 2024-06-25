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

	t, _ := template.New(active.TemplateName).ParseFiles(templateFiles...)

	data["Organisation"] = s.Config.Organisation
	data["PageTitle"] = active.Name + " - "

	slog.Info("static handler",
		slog.String("activeHref", active.Href),
		slog.String("activeTemplate", active.TemplateName),
	)

	s.Write(w, 200, t, active.TemplateName, data)
}
