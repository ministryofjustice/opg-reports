package page

import (
	"bufio"
	"bytes"
	"opg-reports/report/internal/utils"
	"testing"
)

type testPageContent struct {
	Class string
	Name  string
}

// TestHTMLPage uses `test.html` page to render simple dummy data
func TestHTMLPage(t *testing.T) {
	var (
		err          error
		templateDir  = "./testdata"
		byteBuffer   = new(bytes.Buffer)
		buffer       = bufio.NewWriter(byteBuffer)
		templates    = GetTemplateFiles(templateDir)
		templateName = "test"
		data         = &testPageContent{Class: "foobar", Name: "TEST"}
		page         = New(templates, utils.TemplateFunctions())
		expected     = `<h1 class='foobar'>TEST</h1>`
	)
	err = page.WriteToBuffer(buffer, templateName, data)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	buffer.Flush()
	content := byteBuffer.String()
	if content != expected {
		t.Errorf("content does not match\n--> expected:\n%s\n--> actual:\n%s\n", expected, content)
	}
}
