package httpx

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"opg-reports/report/packages/convert"
	"os"
	"path/filepath"
	"testing"
)

var (
	_ Mux = &mux{}
)

type tSimpleStruct struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func setName(ctx context.Context, m Mux, r FitleredRequest, cfg MuxConfig, response *ResponseContent) {
	response.Data["test"] = &tSimpleStruct{Name: "foobar"}
}
func setValue(ctx context.Context, m Mux, r FitleredRequest, cfg MuxConfig, response *ResponseContent) {
	s := response.Data["test"].(*tSimpleStruct)
	s.Value = "added"
}

func htmlTest(ctx context.Context, m Mux, r FitleredRequest, cfg MuxConfig, response *ResponseContent) {
	response.Data["html"] = map[string]string{
		"class": "test-class-name",
		"title": "Page Title!",
	}
}

func TestHttpxMuxRequestHandlerJSON(t *testing.T) {

	// simple test to return a json message object
	mux := NewMux(t.Context(), nil, nil)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/foo/bar/?date_start=2025-12`, nil)
	// attach multiple functions to the endpoint
	mux.Register(`/foo/bar/{$}`, setName, setValue)
	mux.ServeHTTP(rr, req)

	resp := rr.Result()
	// check the result returns correctly
	if resp.StatusCode != http.StatusOK {
		t.Errorf("response code not as expected: [%v]", rr.Code)
	}
	// check content type?
	if resp.Header.Get(contentTypeHeader) != jsonContentType {
		t.Errorf("content type not as expected: [%v]", rr.Header().Get(contentTypeHeader))
	}

	// now check results..
	res := &ResponseContent{}
	tester := &tSimpleStruct{}
	convert.Between(resp, &res)
	convert.Between(res.Data["test"], &tester)

	if tester.Name != "foobar" {
		t.Errorf("returned name does not match expected value")
	}
	if tester.Value != "added" {
		t.Errorf("returned value does not match expected value")
	}

}

func TestHttpxMuxRequestHandlerHTML(t *testing.T) {

	tmpl, e := tmpl(writeTemplate(t.TempDir()))
	if e != nil {
		t.Errorf("unexpected template error: %s", e.Error())
		t.FailNow()
	}

	// simple test to return a json message object
	mux := NewMux(t.Context(), nil, tmpl)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/html/foo/bar/?date_start=2025-12`, nil)
	// attach multiple functions to the endpoint
	mux.Register(`/html/foo/bar/{$}`, htmlTest)
	mux.ServeHTTP(rr, req)

	resp := rr.Result()
	// check the result returns correctly
	if resp.StatusCode != http.StatusOK {
		t.Errorf("response code not as expected: [%v]", rr.Code)
	}
	// check content type?
	if resp.Header.Get(contentTypeHeader) != htmlContentType {
		t.Errorf("content type not as expected: [%v]", rr.Header().Get(contentTypeHeader))
	}

	// now check results..
	// res := &ResponseContent{}
	res := convert.String(resp)
	expected := `<h1 class='h-test-class-name'>Page Title!</h1>`
	if res != expected {
		t.Errorf("html does not match expected value")
	}

}

func tmpl(file string) (*template.Template, error) {
	return template.New("test").Funcs(tmplFunc()).ParseFiles(file)
}

func tmplFunc() (funcs template.FuncMap) {
	funcs = map[string]interface{}{
		// generic value fetcher from the result struct
		// "V": func(key string, name string, data map[string]any) any {
		// 	var v any
		// 	if segment, ok := data[key]; ok {
		// 		var seg = segment.(map[string]interface{})
		// 		if val, ok := seg[name]; ok {
		// 			v = val
		// 		}
		// 	}
		// 	return v
		// },
	}

	return
}

var htmltemplate = `
{{- define "test" -}}
	{{- $pg := index .Data "html" -}}
	{{- $class := index $pg "class" -}}
	{{- $title := index $pg "title" -}}
	<h1 class='h-{{ $class }}'>{{ $title }}</h1>
{{- end -}}
`

func writeTemplate(dir string) (file string) {

	file = filepath.Join(dir, "test.html")
	os.WriteFile(file, []byte(htmltemplate), os.ModePerm)
	return
}
