// Package render used to execute templates for our front end.
//
// Simple struct Render contains the methods used by the front
// end server to take the data and then render the html output.
//
// This does not touch and http response etc
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
	TemplateFiles []string
	Funcs         template.FuncMap
}

// Template returns a template with correct name, functions and partial files ready to be executed
func (self *Render) Template(name string) (tmpl *template.Template, err error) {
	tmpl, err = template.New(name).Funcs(self.Funcs).ParseFiles(self.TemplateFiles...)
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

// New creates a fresh Render wuth the files and funcs passed
func New(files []string, funcs template.FuncMap) *Render {
	return &Render{
		TemplateFiles: files,
		Funcs:         funcs,
	}
}

func named(n string) (name string) {
	name = n
	name = filepath.Base(name)
	name = strings.Split(name, ".")[0]
	return
}
