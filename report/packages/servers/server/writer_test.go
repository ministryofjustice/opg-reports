package server

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/ctxlog"
	"opg-reports/report/packages/types"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	_ types.Templater = &TemplateInfo{}
	_ types.Writer    = &JSONResponseWriter{}
	_ types.Writer    = &HTMLResponseWriter{}
)

type tData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestPackagesServersServerWriterHTMLRepsonse(t *testing.T) {

	var (
		dir        = t.TempDir()
		writer     = httptest.NewRecorder()
		htmlWriter = NewHTMLWriter(ctxlog.New(t.Context(), nil), writer)
		file       = writeTemplate(dir)
	)
	// set templates
	htmlWriter.SetTemplates(&TemplateInfo{Name: "test", FileList: []string{file}})
	// set a bit of data
	htmlWriter.SetData(&tData{ID: 100, Name: "Foobar"})
	// call the response
	htmlWriter.Respond()
	// check out the result
	resp := writer.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("error with response")
	}
	// get content
	content := convert.String(resp)
	if content == "" {
		t.Errorf("unexpected empty content error:")
	}
	// test content
	if !strings.Contains(content, "Foobar") {
		t.Errorf("result should have foobar as h1")
	}
	if !strings.Contains(content, "class='h100'") {
		t.Errorf("result should have class of heading present")
	}

}

func writeTemplate(dir string) (file string) {
	var template = `{{- define "test" -}}<h1 class='h{{ .ID }}'>{{ .Name }}</h1>{{- end -}}`
	file = filepath.Join(dir, "test.html")
	os.WriteFile(file, []byte(template), os.ModePerm)
	return
}

func TestPackagesServersServerWriterJSONRepsonse(t *testing.T) {

	var (
		writer     = httptest.NewRecorder()
		jsonWriter = NewJSONWriter(ctxlog.New(t.Context(), nil), writer)
	)

	jsonWriter.SetData([]*tData{
		{ID: 100, Name: "test a"},
		{ID: 200, Name: "test b"},
	})
	// call the response
	jsonWriter.Respond()
	// check out the result
	res := writer.Result()
	set := []*tData{}
	// convert the reponse
	convert.Between(res, &set)
	// test the results...
	if len(set) != 2 {
		t.Errorf("incorrect amount of data returned")
		t.FailNow()
	}

	if set[0].ID != 100 || set[1].ID != 200 {
		t.Errorf("ordering difference in result")
	}

}
