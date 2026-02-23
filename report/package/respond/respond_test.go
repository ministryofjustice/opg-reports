package respond

import (
	"net/http"
	"net/http/httptest"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/response"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRespondAsHTML(t *testing.T) {

	var (
		err    error
		dir    = t.TempDir()
		ctx    = cntxt.AddLogger(t.Context(), logger.New("error"))
		writer = httptest.NewRecorder()
		req    = httptest.NewRequest(http.MethodGet, "/", nil)
		data   = map[string]string{
			"Name":  "Foobar",
			"Class": "heading",
		}
	)
	file := writeTemplate(dir)

	AsHTML(ctx, req, writer, data, &Args{Template: "test", TemplateFiles: []string{file}})

	resp := writer.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("error with response")
	}

	content, err := response.AsString(resp)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// test content
	if !strings.Contains(content, "Foobar") {
		t.Errorf("result should have foobar as h1")
	}
	if !strings.Contains(content, "class='heading'") {
		t.Errorf("result should have class of heading present")
	}

}

func writeTemplate(dir string) (file string) {
	var template = `{{- define "test" -}}<h1 class='{{ .Class }}'>{{ .Name }}</h1>{{- end -}}`
	file = filepath.Join(dir, "test.html")
	os.WriteFile(file, []byte(template), os.ModePerm)
	return
}

func TestRespondAsJSON(t *testing.T) {

	var (
		err    error
		rec    map[string]string
		ctx    = cntxt.AddLogger(t.Context(), logger.New("error"))
		writer = httptest.NewRecorder()
		req    = httptest.NewRequest(http.MethodGet, "/", nil)
		data   = map[string]string{
			"test": "01",
		}
	)

	AsJSON(ctx, req, writer, data)

	res := writer.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("error with response")
	}

	err = response.As(res, &rec)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	if rec["test"] != data["test"] {
		t.Errorf("data mismtach")
	}

}
