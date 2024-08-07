package main

import (
	"log/slog"
	"net/http"
	"opg-reports/services/api/aws/cost/monthly"
	"opg-reports/services/api/aws/uptime/daily"
	"opg-reports/services/api/github/standards"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/aws/uptime"
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
	var dir string

	// -- aws costs
	dir = "data/aws/cost/monthly/"
	awsCostMonthDir := os.DirFS(dir).(files.IReadFS)
	awsCostMonthlyFs := files.NewFS(awsCostMonthDir, dir)
	awsCostMonthlyStore := data.NewStoreFromFS[*cost.Cost, *files.WriteFS](awsCostMonthlyFs)
	monthly.Register(mux, awsCostMonthlyStore)

	// -- aws uptime
	dir = "data/aws/uptime/daily/"
	awsUptimeDir := os.DirFS(dir).(files.IReadFS)
	awsUptimeFs := files.NewFS(awsUptimeDir, dir)
	awsUptimeStore := data.NewStoreFromFS[*uptime.Uptime, *files.WriteFS](awsUptimeFs)
	daily.Register(mux, awsUptimeStore)

	// -- github standards
	dir = "data/github/standards/"
	ghStandardsDir := os.DirFS(dir).(files.IReadFS)
	ghStandardsFS := files.NewFS(ghStandardsDir, dir)
	ghStandardsStore := data.NewStoreFromFS[*std.Repository, *files.WriteFS](ghStandardsFS)
	standards.Register(mux, ghStandardsStore)

	// -- server
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
