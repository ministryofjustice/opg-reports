package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/config"
)

type registerHandlersFunc func(ctx context.Context, log *slog.Logger, conf *config.Config, info *FrontInfo, mux *http.ServeMux)

func StartServer(
	ctx context.Context,
	log *slog.Logger,
	conf *config.Config,
	info *FrontInfo,
	mux *http.ServeMux,
	server *http.Server,
	registerFuncs ...registerHandlersFunc,
) {
	// call each register function
	var addr = server.Addr
	for _, registerF := range registerFuncs {
		registerF(ctx, log, conf, info, mux)
	}

	log.Info("Starting front server...")
	log.Info(fmt.Sprintf("FRONT: [http://%s/]", addr))

	server.ListenAndServe()
}
