package mw

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/ministryofjustice/opg-reports/shared/logger"
)

func TestMiddlewareDefaultSecurityHeaders(t *testing.T) {
	logger.LogSetup()
	r := httptest.NewRequest(http.MethodGet, "/user/", nil)

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(SecurityHeaders)
	handler.ServeHTTP(rec, r)

	head := rec.Header()

	for header, value := range securityHeaders {
		if !slices.Contains(head[header], value) {
			t.Errorf("error: [%s] not in headers", header)
		}
	}

}
