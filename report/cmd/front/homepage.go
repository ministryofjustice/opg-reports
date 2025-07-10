package main

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
	"opg-reports/report/internal/htmlpage"
	"opg-reports/report/internal/utils"
)

type homepageData struct {
	htmlpage.HtmlPage
}

// handleHomepage
//   - '/'
func handleHomepage(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	info *FrontInfo,
	writer http.ResponseWriter,
	request *http.Request,
) {
	var (
		byteBuffer   = new(bytes.Buffer)
		buffer       = bufio.NewWriter(byteBuffer)
		templates    = htmlpage.GetTemplateFiles(info.TemplateDir)
		templateName = "index"
		data         = htmlpage.DefaultContent(conf, request)
		page         = htmlpage.New(templates, utils.TemplateFunctions())
	)
	data.Teams = info.Teams
	log.Info("processing page", "url", request.URL.String())
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

func RegisterHomepageHandlers(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	info *FrontInfo,
	mux *http.ServeMux,
) {
	log.Info("registering homepage handlers ...")

	// Homepage
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handleHomepage(ctx, log, conf, info, writer, request)
	})

}
