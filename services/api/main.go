package main

import (
	"embed"
	"log/slog"
	"net/http"
	"opg-reports/services/api/aws/cost/monthly"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/env"
	"opg-reports/shared/files"
	"opg-reports/shared/logger"
)

//go:embed data/aws/cost/monthly/*.json
var awsCostMonthlyFs embed.FS

func main() {
	// configure the logger
	logger.LogSetup()

	// mwMux := middleware.NewMiddlewareServerMux(http.NewServeMux())

	mux := http.NewServeMux()

	awsCostMonthlyF := files.NewFS(awsCostMonthlyFs, "data/aws/cost/monthly/")
	awsCostMonthlyS := data.NewStoreFromFS[*cost.Cost, *files.WriteFS](awsCostMonthlyF)
	awsCostMonthlyApi := monthly.New(awsCostMonthlyS, awsCostMonthlyF)
	awsCostMonthlyApi.Register(mux)

	addr := env.Get("ADDR", ":8080")
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	slog.Info("starting api server",
		slog.String("log_level", logger.Level().String()),
		slog.String("address", addr),
	)
	server.ListenAndServe()
}
