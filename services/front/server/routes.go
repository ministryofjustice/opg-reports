package server

import (
	"log/slog"
	"net/http"
	"opg-reports/shared/server"
)

func (s *FrontWebServer) Register(mux *http.ServeMux) {
	all := s.Nav.All()

	for uri, item := range all {
		slog.Debug("registering route", slog.String("uri", uri))
		if item.Api == "" {
			item.Registered = true
			mux.HandleFunc(uri, server.Middleware(s.Static, server.LoggingMW, server.SecurityHeadersMW))
		}
	}
}
