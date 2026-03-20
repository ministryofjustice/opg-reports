package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"opg-reports/report/internal/config"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/fmtx"
	"os"
	"path/filepath"
	"testing"
)

var (
	_ MuxResponder[*testResponse]     = setResponseName
	_ MuxResponder[*testResponse]     = setResponseValue
	_ MuxResponder[*testHTMLResponse] = setResponseHTML
)

type testResponse struct {
	ResponseData
	Name  string `json:"name"`
	Value string `json:"value"`
}

type testHTMLResponse struct {
	ResponseData
	Class string `json:"class"`
	Title string `json:"title"`
}

func (self *testHTMLResponse) TemplateName() string {
	return `test`
}

func setResponseName[T *testResponse](ctx context.Context, cfg MuxConfigurer, r FitleredRequest, response *testResponse) {
	response.Name = "foo"
}
func setResponseValue[T *testResponse](ctx context.Context, cfg MuxConfigurer, r FitleredRequest, response *testResponse) {
	response.Value = "bar"
}
func setResponseHTML[T *testHTMLResponse](ctx context.Context, cfg MuxConfigurer, r FitleredRequest, response *testHTMLResponse) {
	response.Class = `test-class-name`
	response.Title = `Page Title!`
}

func TestHttpxMuxRequestHandlerAsJSON(t *testing.T) {
	// setup the  config & server
	cfg := config.NewApi()
	mux := NewMux()
	// create a test request and recorder to capture it
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/foo/bar/?date_start=2025-12`, nil)
	// register the end point - chain functions to update different parts
	Register(t.Context(), mux, cfg, `/foo/bar/{$}`, nil, setResponseName, setResponseValue)
	// process the end point
	mux.ServeHTTP(rr, req)
	// grab the result
	resp := rr.Result()
	// check the result returns correctly
	if resp.StatusCode != http.StatusOK {
		t.Errorf("response code not as expected: [%v]", rr.Code)
	}
	// check content type?
	if resp.Header.Get(contentTypeHeader) != jsonContentType {
		t.Errorf("content type not as expected: [%v]", rr.Header().Get(contentTypeHeader))
	}

	res := &testResponse{}
	convert.Between(resp, &res)

	if res.Name != "foo" {
		t.Errorf("returned name does not match expected value")
	}
	if res.Value != "bar" {
		t.Errorf("returned value does not match expected value")
	}

}

func TestHttpxMuxRequestHandlerAsHTML(t *testing.T) {
	cfg := config.NewFront()
	// write the template and change the directory path
	dir := t.TempDir()
	writeTemplate(dir)
	cfg.Root = `/`
	cfg.TemplateAssetPath = dir
	// create a new mux
	mux := NewMux()
	// create a test request and recorder to capture it
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/html/foo/bar/?date_start=2025-12`, nil)

	fmtx.Dump(cfg)
	// register the end point
	Register(t.Context(), mux, cfg, `/html/foo/bar/{$}`, nil, setResponseHTML)
	// process the end point
	mux.ServeHTTP(rr, req)
	// grab the result
	resp := rr.Result()
	// check the result returns correctly
	if resp.StatusCode != http.StatusOK {
		t.Errorf("response code not as expected: [%v]", rr.Code)
	}
	// check content type?
	if resp.Header.Get(contentTypeHeader) != htmlContentType {
		t.Errorf("content type not as expected: [%v]", rr.Header().Get(contentTypeHeader))
	}
	res := convert.String(resp)
	expected := `<h1 class='h-test-class-name'>Page Title!</h1>`
	if res != expected {
		t.Errorf("html does not match expected value: [%s]", res)
	}
}

var htmltemplate = `
{{- define "test" -}}
	<h1 class='h-{{ .Class }}'>{{ .Title }}</h1>
{{- end -}}
`

func writeTemplate(dir string) (file string) {
	file = filepath.Join(dir, "test.html")
	os.WriteFile(file, []byte(htmltemplate), os.ModePerm)
	return
}
