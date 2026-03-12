package httpx

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types"
	"os"
	"path/filepath"
	"testing"
)

var (
	w                           = &ResponseWriter{}
	_ http.ResponseWriter       = w
	_ types.HttpxResponseWriter = w
)

type testItem struct {
	Name string `json:"name"`
}

type testTmpl struct {
	name  string
	files []string
}

func (self *testTmpl) Name() string {
	return self.name
}
func (self *testTmpl) WithName(n string) types.HttpxTemplater {
	self.name = n
	return self
}
func (self *testTmpl) Template() (tp *template.Template) {
	tp, _ = template.New(self.name).Funcs(nil).ParseFiles(self.files...)
	return
}

func TestHttpxResponseWriterHTML(t *testing.T) {
	dir := t.TempDir()
	file := writeTemplate(dir)

	fx := func(ctx types.Contextx, t types.HttpxTemplater, w types.HttpxResponseWriter, r types.HttpxRequest) {
		var data = map[string]string{
			"Name":  "Foobar",
			"Class": "heading",
		}
		w.Send(data, t)
	}
	tp := &testTmpl{name: "test", files: []string{file}}
	mux := NewServeMux()
	mux.HandleFuncx("/foo/{$}", tp, fx)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/foo/", nil)
	mux.ServeHTTP(rr, req)

	str := convert.String(rr)

	if str != `<h1 class='heading'>Foobar</h1>` {
		t.Errorf("html not returned correctly:\n%s", str)
	}

}

func TestHttpxResponseWriterJSON(t *testing.T) {
	mux := NewServeMux()
	// add a dummy func that uses send..
	fx := func(ctx types.Contextx, t types.HttpxTemplater, w types.HttpxResponseWriter, r types.HttpxRequest) {
		var data = []testItem{}
		for i := 0; i < 100000; i++ {
			data = append(data, testItem{Name: fmt.Sprintf("name-%d", 100000+i)})
		}
		w.Send(data, t)
	}
	mux.HandleFuncx("/foo/{$}", nil, fx)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/foo/", nil)
	mux.ServeHTTP(rr, req)

	items := []testItem{}
	convert.Between(rr, &items)

	if len(items) != 100000 {
		t.Errorf("incorrect number of items returned")
	}

}

func writeTemplate(dir string) (file string) {
	var template = `{{- define "test2" -}}Should not be rendered{{- end -}}{{- define "test" -}}<h1 class='{{ .Class }}'>{{ .Name }}</h1>{{- end -}}`
	file = filepath.Join(dir, "test.html")
	os.WriteFile(file, []byte(template), os.ModePerm)
	return
}
