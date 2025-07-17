package page

import (
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// Page is used to deal with rendering the go html/template stack
// with local template files and combining data into the page.
//
// WriteToBuffer pushes the parsed template content into a buffer which
// can then be used in http.ResponseWriter.Write to include in the front
// end response
type Page struct {
	files []string
	funcs template.FuncMap
}

// GetTemplate loads all template files and function names into the template stack
// and will generate a template.Template based on the name provided which should
// be the base name and match the defined label
func (self *Page) GetTemplate(name string) (tmpl *template.Template, err error) {
	tmpl, err = template.New(name).
		Funcs(self.funcs).
		ParseFiles(self.files...)
	return
}

// WriteToBuffer uses the template name and data to render the html template stack and
// push the resulting content of that into the buffer passed.
//
// The buffer is then used in http.ResponseWriter.Write to render a response for the
// current request
//
// If an error occurs with finding the template, not template is generated
func (self *Page) WriteToBuffer(buffer io.Writer, templateName string, data any) (err error) {
	var templ *template.Template

	templateName = pageName(templateName)
	templ, err = self.GetTemplate(templateName)
	if err != nil {
		return
	}
	err = templ.ExecuteTemplate(buffer, templateName, data)
	return
}

func New(templateFiles []string, templateFunctions template.FuncMap) *Page {
	return &Page{
		files: templateFiles,
		funcs: templateFunctions,
	}
}

// pageName helper function to ensure the template name hasn't
// got any file extentions of /'s etc
func pageName(n string) (name string) {
	name = n
	name = filepath.Base(name)
	name = strings.Split(name, ".")[0]
	return
}
