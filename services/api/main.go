package main

import (
	"embed"
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
	"opg-reports/shared/server/response"
)

//go:embed data/aws/cost/monthly/*.json
var awsCostMonthlyFs embed.FS

//go:embed data/github/standards/*.json
var ghStandardsFs embed.FS

func main() {
	// configure the logger
	logger.LogSetup()

	mux := http.NewServeMux()

	awsCostMonthlyFs := files.NewFS(awsCostMonthlyFs, "data/aws/cost/monthly/")
	awsCostMonthlyStore := data.NewStoreFromFS[*cost.Cost, *files.WriteFS](awsCostMonthlyFs)
	awsResp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	awsCostMonthlyApi := monthly.New(awsCostMonthlyStore, awsCostMonthlyFs, awsResp)
	awsCostMonthlyApi.Register(mux)

	ghStandardsFS := files.NewFS(ghStandardsFs, "data/github/standards/")
	ghStandardsStore := data.NewStoreFromFS[*std.Repository, *files.WriteFS](ghStandardsFS)
	ghResp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	ghStandardsApi := standards.New(ghStandardsStore, ghStandardsFS, ghResp)
	ghStandardsApi.Register(mux)

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
