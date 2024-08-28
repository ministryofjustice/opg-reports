package httphandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

func TestSharedSrvrHttphandler(t *testing.T) {
	mock := mockOKResponse()

	handler := httphandler.New("", "", mock.URL)
	err := handler.Get()
	if err != nil {
		t.Errorf("failed to get http data [%s]", err.Error())
	}

	str, _ := convert.Stringify(handler.Response)
	if str != "OK" {
		t.Errorf("content did not match: [%s]", str)
	}
	if handler.StatusCode != http.StatusOK {
		t.Errorf("status code mismatch")
	}

}

func mockOKResponse() *httptest.Server {
	return testhelpers.MockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}, "warn")
}
