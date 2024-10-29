package render

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

type Render struct {
	PartialFiles []string
	Funcs        template.FuncMap
}

// Template returns a template with correct name, functions and partial files ready to be executed
func (self *Render) Template(name string) (tmpl *template.Template, err error) {
	tmpl, err = template.New(name).Funcs(self.Funcs).ParseFiles(self.PartialFiles...)
	return
}

// Write generates the content of the partial and writes that into the writer passed
// making use of the data to update ant content
func (self *Render) Write(name string, data any, writer io.Writer) (err error) {
	var nm = named(name)
	var tmpl *template.Template

	tmpl, err = self.Template(nm)
	if err != nil {
		return
	}

	err = tmpl.ExecuteTemplate(writer, nm, data)

	return
}

// Execute generates the content of the partial directly and returns as a string
// making use of bufio to generate a localised buffer
// In most cases, should utilise Write instead that accepts a buffer
func (self *Render) Execute(name string, data any) (str string, err error) {
	var nm = named(name)
	var buf = new(bytes.Buffer)
	var writer = bufio.NewWriter(buf)

	self.Write(nm, data, writer)
	// make sure to flush the writer before grabbing output
	writer.Flush()
	// get output of the buffer
	str = buf.String()
	return
}

func NewPartial(partials []string, funcs template.FuncMap) *Render {
	return &Render{
		PartialFiles: partials,
		Funcs:        funcs,
	}
}

func named(n string) (name string) {
	name = n
	name = filepath.Base(name)
	name = strings.Split(name, ".")[0]
	return
}
