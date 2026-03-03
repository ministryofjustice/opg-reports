package main

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/respond"
)

type PingResponse struct {
	Version string
	SHA     string
}

func registerPingAndHome(ctx context.Context, mux *http.ServeMux, config *cli) {
	var log = cntxt.GetLogger(ctx)

	log.Info(fmt.Sprintf("[%s] registering home [%s] to handler", "main", "/"))
	mux.HandleFunc("/{$}", func(writer http.ResponseWriter, request *http.Request) {
		pingResponse(ctx, config, request, writer)
	})
	log.Info(fmt.Sprintf("[%s] registering ping [%s] to handler", "main", "/ping/"))
	mux.HandleFunc("/ping/{$}", func(writer http.ResponseWriter, request *http.Request) {
		pingResponse(ctx, config, request, writer)
	})
}

func pingResponse(ctx context.Context, conf *cli, request *http.Request, writer http.ResponseWriter) {
	result := &PingResponse{
		Version: conf.Version,
		SHA:     conf.SHA,
	}
	respond.AsJSON(ctx, request, writer, result)
}
