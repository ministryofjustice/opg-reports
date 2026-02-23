package respond

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
)

type Args struct {
	Template      string
	TemplateFiles []string
	Funcs         template.FuncMap
}

// AsHTML will execute template and return as html cotnent using data passed along
func AsHTML(ctx context.Context, request *http.Request, writer http.ResponseWriter, data any, args *Args) {

	var (
		err        error
		dataBytes  []byte
		tmpl       *template.Template
		buffer     = new(bytes.Buffer)
		buffWriter = bufio.NewWriter(buffer)
	)
	// parse the template
	tmpl, err = template.New(args.Template).Funcs(args.Funcs).ParseFiles(args.TemplateFiles...)
	if err != nil {
		return
	}
	// execute it
	err = tmpl.ExecuteTemplate(buffer, args.Template, data)

	if err != nil {
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

// AsJSON writes data as pure json, no html / template stack involved
func AsJSON(ctx context.Context, request *http.Request, writer http.ResponseWriter, data any) {
	var (
		err        error
		dataBytes  []byte
		buffer     = new(bytes.Buffer)
		buffWriter = bufio.NewWriter(buffer)
	)

	// convert the data struct into jsonified bytes
	dataBytes, err = json.MarshalIndent(data, "", "  ")

	if err != nil {
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
