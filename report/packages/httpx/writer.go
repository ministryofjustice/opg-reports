package httpx

import (
	"bufio"
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
)

const (
	contentTypeHeader string = `Content-Type`
	jsonContentType   string = `application/json`
	htmlContentType   string = `text/html`
)

// ResponseWriter is an extension of the http.ResponseWriter
// that provides methods to allow for easier switching
// between json & html response and reducing the code within
// the main handler.
type ResponseWriter interface {
	http.ResponseWriter
	// ContentType returns the value to use
	// for the content-type header
	ContentType() (string, string)
	// Status returns the status code to use
	// for the response
	Status() int
	// Bytes converts the data into []bytes
	// which can then be passed along to Write
	Bytes(sourceData any) ([]byte, error)
	// BytesAndHeaders converts the data into []bytes and
	// also sets the headers
	BytesAndHeaders(sourceData any) (out []byte, err error)
}

// errorResponse is a simple struct used to return
// error within the writer process so theres no
// type switching.
// used by the json resposne writter
type errorResponse struct {
	Error string `json:"error"`
}

// jsonResponse is the JSON response handler that will attempt
// to return an data as json content
type jsonResponse struct {
	http.ResponseWriter
	code        int    // http status code, defaults to http.StatusOK
	contentType string // http content type - always json
}

// ContentType returns the nmae & value to use
// for the content-type header
func (self *jsonResponse) ContentType() (string, string) {
	return contentTypeHeader, self.contentType
}

// Status returns the status code to use
// for the response
func (self *jsonResponse) Status() int {
	return self.code
}

// Bytes converts the data into []bytes
// which can then be passed along to Write
func (self *jsonResponse) Bytes(sourceData any) (out []byte, err error) {
	var (
		content      []byte
		buffer       = new(bytes.Buffer)
		bufferWriter = bufio.NewWriter(buffer)
	)
	content, err = json.MarshalIndent(sourceData, "", "  ")
	// if theres an error then use error response struct to capture it
	// and set that there has been an error
	if err != nil {
		self.code = http.StatusInternalServerError
		content, _ = json.MarshalIndent(&errorResponse{Error: err.Error()}, "", "  ")
	}
	// write to the buffer
	bufferWriter.Write(content)
	bufferWriter.Flush()
	// return the buffer content
	out = buffer.Bytes()
	return
}

// BytesAndHeaders
func (self *jsonResponse) BytesAndHeaders(sourceData any) (out []byte, err error) {
	out, err = self.Bytes(sourceData)
	self.Header().Set(self.ContentType())
	self.WriteHeader(self.Status())
	return
}

// htmlResponse expands on jsonResposne to add html template parsing
type htmlResponse struct {
	*jsonResponse
	template *template.Template // template is a pre-complied template
}

// Bytes converts the data into []bytes
// which can then be passed along to Write
func (self *htmlResponse) Bytes(sourceData any) (out []byte, err error) {
	var (
		buffer       = new(bytes.Buffer)
		bufferWriter = bufio.NewWriter(buffer)
	)
	err = self.template.ExecuteTemplate(bufferWriter, self.template.Name(), sourceData)
	// if theres an error then write that error to the buffer
	if err != nil {
		self.code = http.StatusInternalServerError
		bufferWriter.WriteString(err.Error())
	}
	// flush the buffer
	bufferWriter.Flush()
	// return the buffer content
	out = buffer.Bytes()
	return
}

// BytesAndHeaders
func (self *htmlResponse) BytesAndHeaders(sourceData any) (out []byte, err error) {
	out, err = self.Bytes(sourceData)
	self.Header().Set(self.ContentType())
	self.WriteHeader(self.Status())
	return
}

// NewResponseWriter creates a new response writer from the w and the template details
func NewResponseWriter(w http.ResponseWriter, tmpl *template.Template) ResponseWriter {
	if tmpl != nil {
		return NewHTMLResponseWriter(w, tmpl)
	}
	return NewJSONResponseWriter(w)
}

// NewJSONResponseWriter returns a json response writer that will
// return all data as a json output
func NewJSONResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &jsonResponse{
		ResponseWriter: w,
		code:           http.StatusOK,
		contentType:    jsonContentType,
	}
}

// NewHTMLResponseWriter returns a html response writer
func NewHTMLResponseWriter(w http.ResponseWriter, tmpl *template.Template) ResponseWriter {
	return &htmlResponse{
		template: tmpl,
		jsonResponse: &jsonResponse{
			ResponseWriter: w,
			code:           http.StatusOK,
			contentType:    htmlContentType,
		},
	}
}
