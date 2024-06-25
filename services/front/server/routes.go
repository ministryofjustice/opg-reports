package server

import (
	"log/slog"
	"net/http"
	"opg-reports/shared/server"
	"strings"
)

func (s *FrontWebServer) Register(mux *http.ServeMux) {
	all := s.Nav.All()

	for uri, item := range all {
		// trim & readd the ending parts of the url pattern
		uri = strings.TrimSuffix(strings.TrimSuffix(uri, "{$}"), "/")
		uri += "/{$}"
		item.Registered = true
		hf := s.Static
		if item.Api != "" {
			hf = s.Dynamic
		}
		slog.Info("registering route", slog.String("uri", uri), slog.Any("api", item.Api))
		mux.HandleFunc(uri, server.Middleware(hf, server.LoggingMW, server.SecurityHeadersMW))

	}
}
