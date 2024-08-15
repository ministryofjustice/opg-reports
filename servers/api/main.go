package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/exists"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

var databases map[string]string = map[string]string{
	"github_standards": "./github_standards.db",
	"aws_costs":        "./aws_costs.db",
}

func main() {
	logger.LogSetup()
	ctx := context.Background()
	slog.Info("databases", slog.String("db:", fmt.Sprintf("%+v", databases)))

	mux := http.NewServeMux()
	// -- github standards
	if !exists.FileOrDir(databases["github_standards"]) {
		slog.Error("database missing for github_standards", slog.String("db", databases["github_standards"]))
		os.Exit(1)
	}
	github_standards.Register(ctx, mux, databases["github_standards"])
	// -- aws costs
	if !exists.FileOrDir(databases["aws_costs"]) {
		slog.Error("database missing for aws_costs", slog.String("db", databases["aws_costs"]))
		os.Exit(1)
	}
	aws_costs.Register(ctx, mux, databases["aws_costs"])

	// -- start the server
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
