package main

import (
	"bufio"
	"bytes"
	"net/http"
	"opg-reports/report/internal/page"
	"opg-reports/report/internal/utils"
)

// Respond handles running the html/template stack with all
// the functions and data and writes the result to the response
// writers buffer for returning to user
func Respond(
	writer http.ResponseWriter,
	request *http.Request,
	templateName string,
	templates []string,
	data any,
) {
	var (
		byteBuffer = new(bytes.Buffer)
		buffer     = bufio.NewWriter(byteBuffer)
		page       = page.New(templates, utils.TemplateFunctions())
	)
	// call page WriteToBuffer to run the template stack and write to buffer
	page.WriteToBuffer(buffer, templateName, data)
	// write ok status & content type to response
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	// force flush the underlying buffer to make sure all content is updated
	buffer.Flush()
	// write content to the response
	writer.Write(byteBuffer.Bytes())
}
