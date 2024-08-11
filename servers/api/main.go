package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const github_standards_db = "./dbs/github_standards.db"

func main() {
	logger.LogSetup()
	ctx := context.Background()

	mux := http.NewServeMux()
	github_standards.Register(ctx, mux, github_standards_db)

	addr := env.Get("API_ADDR", consts.API_ADDR)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	slog.Info("starting api server",
		slog.String("log_level", logger.Level().String()),
		slog.String("api_address", addr),
	)
	server.ListenAndServe()
}
