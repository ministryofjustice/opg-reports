package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/seeder/github_standards_seed"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/exists"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const (
	github_standards_db     string = "github_standards.db"
	github_standards_schema string = "github_standards.sql"
	github_standards_N      int    = 1500
)

func init() {
	// -- seed databases
	if !exists.FileOrDir(github_standards_db) && exists.FileOrDir(github_standards_schema) {
		slog.Info("creating a seeded database...")
		db, err := github_standards_seed.NewSeed("./", github_standards_N)
		defer db.Close()
		if err != nil {
			slog.Error("error with seeding:" + err.Error())
		}
	}

}

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
