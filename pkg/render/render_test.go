package render_test

import (
	"bufio"
	"bytes"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/pkg/render"
	"github.com/ministryofjustice/opg-reports/pkg/tmplfuncs"
)

const tmpDir string = "../../servers/sfront/templates/partials"

// TestRenderRender checks that the test partial
// is rendered correctly as a string
func TestRenderExecute(t *testing.T) {
	pattern := filepath.Join(tmpDir, "*.gotmpl")
	files, _ := filepath.Glob(pattern)
	funcs := tmplfuncs.All
	dummy := map[string]string{"Name": "TEST", "Class": "foobar"}
	expected := `<h1 class='foobar'>TEST</h1>`

	var p = "test"
	rnd := render.New(files, funcs)

	actual, err := rnd.Execute(p, dummy)

	if err != nil {
		t.Errorf("error with template: [%s]", err.Error())
	}

	if expected != actual {
		t.Errorf("template render failed:\n expected:\n  -%s-\n actual:\n -%s-", expected, actual)
	}
}

func TestRenderWrite(t *testing.T) {
	pattern := filepath.Join(tmpDir, "*.gotmpl")
	files, _ := filepath.Glob(pattern)
	funcs := tmplfuncs.All
	dummy := map[string]string{"Name": "TEST", "Class": "foobar"}
	expected := `<h1 class='foobar'>TEST</h1>`
	buf := new(bytes.Buffer)
	wr := bufio.NewWriter(buf)

	var p = "test"
	rnd := render.New(files, funcs)

	err := rnd.Write(p, dummy, wr)
	wr.Flush()
	actual := buf.String()

	if err != nil {
		t.Errorf("error with template: [%s]", err.Error())
	}

	if expected != actual {
		t.Errorf("template render failed:\n expected:\n  -%s-\n actual:\n -%s-", expected, actual)
	}
}
