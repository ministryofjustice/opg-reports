package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/api/aws_uptime"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/exists"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

var databases map[string]string = map[string]string{
	"github_standards": "./github_standards.db",
	"aws_costs":        "./aws_costs.db",
	"aws_uptime":       "./aws_uptime.db",
}

func main() {
	logger.LogSetup()
	ctx := context.Background()
	slog.Info("databases", slog.String("db:", fmt.Sprintf("%+v", databases)))

	mux := http.NewServeMux()

	// -- github standards
	if !exists.FileOrDir(databases["github_standards"]) {
		slog.Error("database missing for github_standards", slog.String("db", databases["github_standards"]))
	} else {
		ghsServer := api.New(ctx, databases["github_standards"])
		github_standards.Register(mux, ghsServer)
	}

	// -- aws costs
	if !exists.FileOrDir(databases["aws_costs"]) {
		slog.Error("database missing for aws_costs", slog.String("db", databases["aws_costs"]))
	} else {
		awscServer := api.New(ctx, databases["aws_costs"])
		aws_costs.Register(mux, awscServer)
	}

	// -- aws uptime
	if !exists.FileOrDir(databases["aws_uptime"]) {
		slog.Error("database missing for aws_uptime", slog.String("db", databases["aws_uptime"]))
	} else {
		awsuServer := api.New(ctx, databases["aws_uptime"])
		aws_uptime.Register(mux, awsuServer)
	}

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
