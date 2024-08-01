package server

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/shared/logger"
	"slices"
	"testing"
)

func TestMiddlewareDefaultSecurityHeaders(t *testing.T) {
	logger.LogSetup()
	r := httptest.NewRequest(http.MethodGet, "/user/", nil)

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(SecurityHeadersMW)
	handler.ServeHTTP(rec, r)

	head := rec.Header()

	for header, value := range securityHeaders {
		if !slices.Contains(head[header], value) {
			t.Errorf("error: [%s] not in headers", header)
		}
	}

}
