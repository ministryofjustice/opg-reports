package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
)

type registerHandlersFunc func(ctx context.Context, log *slog.Logger, conf *config.Config, mux *http.ServeMux)

func StartServer(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	mux *http.ServeMux,
	server *http.Server,
	registerFuncs ...registerHandlersFunc,
) {
	// call each register function
	var addr = server.Addr
	for _, registerF := range registerFuncs {
		registerF(ctx, log, conf, mux)
	}

	log.Info("Starting front server...")
	log.Info(fmt.Sprintf("ROOT ASSET DIR: [%s]", assetRoot))
	log.Info(fmt.Sprintf("GOVUK ASSET DIR: [%s]", govUKAssetDir))
	log.Info(fmt.Sprintf("LOCAL ASSET DIR: [%s]", localAssetsDir))
	log.Info(fmt.Sprintf("TEMPLATE DIR: [%s]", templateDir))
	log.Info(fmt.Sprintf("FRONT: [http://%s/]", addr))

	server.ListenAndServe()
}
