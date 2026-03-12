package httpx

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"opg-reports/report/packages/types"
)

// ResponseWriter is an extended writer that returns
// either json or html depending on the existence of the
// template value being passed to send.
//
// When using this, call `Send` instead of write.
//
// TODO:
//   - add context with logging
type ResponseWriter struct {
	http.ResponseWriter
}

// Send converts data into html / json response and then calls write
func (self *ResponseWriter) Send(data any, template types.HttpxTemplater) (code int, err error) {
	var (
		// 	content     []byte
		contentType  string
		writer       = self.ResponseWriter
		buffer       = new(bytes.Buffer)
		bufferWriter = bufio.NewWriter(buffer)
	)
	code = http.StatusOK

	if template != nil {
		contentType, err = asHTML(data, bufferWriter, template)
	} else {
		contentType, err = asJSON(data, bufferWriter)
	}

	if err != nil {
		code = http.StatusInternalServerError
	}

	writer.WriteHeader(code)
	writer.Header().Set("Content-Type", contentType)

	bufferWriter.Flush()
	code, err = writer.Write(buffer.Bytes())
	return
}

// NewResponseWriter returns a new instance of ResponseWriter
func NewResponseWriter(w http.ResponseWriter) types.HttpxResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
	}
}

// asHTML tries to execute the template stack passed along and writes that to the buffer
//
// On failure, writes the error message instead.
func asHTML(data any, bufferWriter *bufio.Writer, template types.HttpxTemplater) (contentType string, err error) {
	var tmpl = template.Template()
	var name = template.Name()

	contentType = "text/html"
	err = tmpl.ExecuteTemplate(bufferWriter, name, data)
	if err != nil {
		bufferWriter.WriteString(err.Error())
	}
	return
}

// asJSON tries to convert the data into a json result and will write that to the buffer.
//
// On failure, writes the error message instead and adjusts the content type to be html
func asJSON(data any, bufferWriter *bufio.Writer) (contentType string, err error) {
	var content []byte
	contentType = "application/json"
	content, err = json.MarshalIndent(data, "", "  ")

	if err != nil {
		bufferWriter.WriteString(err.Error())
		contentType = "text/html"
	} else {
		bufferWriter.Write(content)
	}
	return
}
