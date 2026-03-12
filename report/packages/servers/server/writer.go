package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"opg-reports/report/packages/types"
	"text/template"
)

// TemplateInfo provides details about html templates
// used when writing the result
type TemplateInfo struct {
	Name      string
	FileList  []string
	Functions template.FuncMap
}

// Template is the defined name (`{{ define "$x" }}`) that we want to render.
// Template fiels can contain many templates.
func (self *TemplateInfo) Template() string {
	return self.Name
}

// Files is a list of all template fiels in the stack to render
func (self *TemplateInfo) Files() []string {
	return self.FileList
}

// Func returns the list of functions to setup
func (self *TemplateInfo) Funcs() template.FuncMap {
	return self.Functions
}

// HTMLResponseWriter handles api style data results, so those
// that are pure json
type HTMLResponseWriter struct {
	http.ResponseWriter
	// context
	ctx types.ContextLogger
	// the raw data that will get marshaled and returned
	data any
	// template data
	templates types.Templater
}

func (self *HTMLResponseWriter) Ctx() types.ContextLogger {
	return self.ctx
}

// SetData updates the data property
func (self *HTMLResponseWriter) SetData(data any) {
	self.data = data
}

// SetTemplates does nothing for json response writer
func (self *HTMLResponseWriter) SetTemplates(templates types.Templater) {
	self.templates = templates
}

// Respond handles writing the currently set data to the writer as
// the response
func (self *HTMLResponseWriter) Respond() {
	var (
		err        error
		dataBytes  []byte
		log        = self.ctx.Log()
		args       = self.templates
		tmpl       *template.Template
		buffer     = new(bytes.Buffer)
		buffWriter = bufio.NewWriter(buffer)
		data       = self.data
		writer     = self.ResponseWriter
	)
	// parse the template
	tmpl, err = template.
		New(args.Template()).
		Funcs(args.Funcs()).
		ParseFiles(args.Files()...)

	// return error
	if err != nil {
		log.Error("HTMLResponseWriter: error parsing templates", "err", err.Error())
		return
	}
	// exec the template
	err = tmpl.ExecuteTemplate(buffer, args.Template(), data)
	if err != nil {
		log.Error("HTMLResponseWriter: error executing templates", "err", err.Error())
		buffWriter.WriteString(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "text/html")
	} else {
		buffWriter.Write(dataBytes)
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
	}
	buffWriter.Flush()
	writer.Write(buffer.Bytes())
}

// JSONResponseWriter handles api style data results, so those
// that are pure json
type JSONResponseWriter struct {
	http.ResponseWriter
	// context
	ctx types.ContextLogger
	// the raw data that will get marshaled and returned
	data any
}

func (self *JSONResponseWriter) Ctx() types.ContextLogger {
	return self.ctx
}

// SetData updates the data property
func (self *JSONResponseWriter) SetData(data any) {
	self.data = data
}

// SetTemplates does nothing for json response writer
func (self *JSONResponseWriter) SetTemplates(cfg types.Templater) {}

// Respond handles writing the currently set data to the writer as
// the response
func (self *JSONResponseWriter) Respond() {
	var (
		err        error
		dataBytes  []byte
		log        = self.ctx.Log()
		buffer     = new(bytes.Buffer)
		buffWriter = bufio.NewWriter(buffer)
		data       = self.data
		writer     = self.ResponseWriter
	)

	// convert the data struct into jsonified bytes
	dataBytes, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error("JSONResponseWriter: error marshaling data", "err", err.Error())
		buffWriter.WriteString(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "text/html")
	} else {
		buffWriter.Write(dataBytes)
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
	}
	buffWriter.Flush()
	writer.Write(buffer.Bytes())
}

func NewJSONWriter(ctx types.ContextLogger, w http.ResponseWriter) types.Writer {
	return &JSONResponseWriter{
		ResponseWriter: w,
		ctx:            ctx,
	}
}

func NewHTMLWriter(ctx types.ContextLogger, w http.ResponseWriter) types.Writer {
	return &HTMLResponseWriter{
		ResponseWriter: w,
		ctx:            ctx,
	}
}
