package httpx

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types"
	"testing"
)

var (
	s                = &ServeMux{}
	_ http.Handler   = s
	_ types.HttpxMux = s
)

type testTpl struct{}

func (self *testTpl) Template() *template.Template {
	return nil
}

func TestHttpxServeMuxHandleFuncX(t *testing.T) {

	// add a dummy func (types.HttpxHandlerFunc)
	fx := func(ctx types.Contextx, t types.HttpxTemplater, w types.HttpxResponseWriter, r types.HttpxRequest) {
		w.Write([]byte("foobar"))
	}
	mux := NewServeMux()
	mux.HandleFuncx("/foo/{$}", nil, fx)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/foo/", nil)
	mux.ServeHTTP(rr, req)

	str := convert.String(rr)
	if str != "foobar" {
		t.Errorf("incorrect response")
	}

}

func TestHttpxServeMuxHandleFunc(t *testing.T) {

	mux := NewServeMux()
	// add a dummy func
	mux.HandleFunc("/foo/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bar"))
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/foo/", nil)
	mux.ServeHTTP(rr, req)

	str := convert.String(rr)
	if str != "bar" {
		t.Errorf("incorrect response")
	}

}
