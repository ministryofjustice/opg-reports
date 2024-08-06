package main

import (
	"log/slog"
	"net/http"
	"opg-reports/services/api/aws/cost/monthly"
	"opg-reports/services/api/github/standards"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/env"
	"opg-reports/shared/files"
	"opg-reports/shared/github/std"
	"opg-reports/shared/logger"
	"os"
)

func main() {
	// configure the logger
	logger.LogSetup()

	mux := http.NewServeMux()

	awsCostMonthDir := os.DirFS("data/aws/cost/monthly/").(files.IReadFS)
	awsCostMonthlyFs := files.NewFS(awsCostMonthDir, "data/aws/cost/monthly/")
	awsCostMonthlyStore := data.NewStoreFromFS[*cost.Cost, *files.WriteFS](awsCostMonthlyFs)
	monthly.Register(mux, awsCostMonthlyStore)

	ghStandardsDir := os.DirFS("data/github/standards/").(files.IReadFS)
	ghStandardsFS := files.NewFS(ghStandardsDir, "data/github/standards/")
	ghStandardsStore := data.NewStoreFromFS[*std.Repository, *files.WriteFS](ghStandardsFS)
	standards.Register(mux, ghStandardsStore)

	addr := env.Get("API_ADDR", ":8081")
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
