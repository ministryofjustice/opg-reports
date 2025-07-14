package main

import (
	"bufio"
	"bytes"
	"net/http"
	"opg-reports/report/internal/page"
	"opg-reports/report/internal/utils"
)

func pgSetup(templates []string) (byteBuffer *bytes.Buffer, buffer *bufio.Writer, pg *page.Page) {
	byteBuffer = new(bytes.Buffer)
	buffer = bufio.NewWriter(byteBuffer)
	pg = page.New(templates, utils.TemplateFunctions())
	return

}

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
	var byteBuffer, buffer, pg = pgSetup(templates)
	// call page WriteToBuffer to run the template stack and write to buffer
	err := pg.WriteToBuffer(buffer, templateName, data)
	// If there are no errors rendering the template name and data stack, then
	//   - write the header status as ok
	//   - set the header type
	// 	 - flush the buffer to make sure content is pushed
	// 	 - write templated content
	// If there is an error, we do similar, but force a new template "error"
	// that is empty and should work
	if err == nil {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "text/html")
		buffer.Flush()
		writer.Write(byteBuffer.Bytes())

	} else {
		pg.WriteToBuffer(buffer, "error", data)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "text/html")
		buffer.Flush()
		writer.Write(byteBuffer.Bytes())
	}
}
