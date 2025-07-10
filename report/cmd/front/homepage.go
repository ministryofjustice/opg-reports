package main

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
)

func RegisterHomepageHandlers(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	mux *http.ServeMux,
) {
	log.Info("registering homepage handlers ...")

	// Homepage
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var (
			buf = new(bytes.Buffer)
			wr  = bufio.NewWriter(buf)
		)
		// write content to the buffer
		wr.WriteString("<p>Hello</p>")
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "text/html")

		wr.Flush()
		writer.Write(buf.Bytes())

	})

}
