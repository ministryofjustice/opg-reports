package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/seeder"
	"github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/exists"
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

var databases map[string]string = map[string]string{
	"github_standards": "./github_standards/github_standards.db",
}

func init() {
	logger.LogSetup()
	var lines map[string][]string = map[string][]string{
		"github_standards": {},
	}

	// -- seed databases
	slog.Info("api init")

	// -- create some dummy data for each type
	// --- github standards
	owner := fake.String(12)
	for x := 0; x < 1500; x++ {
		id := 1000 + x
		g := ghs.Fake(&id, &owner)
		if x == 0 {
			lines["github_standards"] = append(lines["github_standards"], g.CSVHead())
		}
		lines["github_standards"] = append(lines["github_standards"], g.ToCSV())

	}

	// -- list of what seeds to create
	var seedList []*seeder.Seed = []*seeder.Seed{
		{
			Label:  "built",
			Table:  "github_standards",
			DB:     databases["github_standards"],
			Schema: "./github_standards/github_standards.sql",
			Source: []string{"./github_standards/github_standards.csv"},
			Dummy:  []string{},
		},
		{
			Label:  "local",
			Table:  "github_standards",
			DB:     "./github_standards.db",
			Schema: "../../datastore/github_standards/github_standards.sql",
			Source: []string{},
			Dummy:  lines["github_standards"],
		},
	}

	for _, sl := range seedList {
		slog.Debug("seed", slog.String("group", sl.Table), slog.String("label", sl.Label))

		// if the schema exists, but the db doesn't, then we create it
		if exists.FileOrDir(sl.Schema) && !exists.FileOrDir(sl.DB) {
			slog.Info("generating seed", slog.String("group", sl.Table), slog.String("label", sl.Label))
			db, err := seeder.New(sl)
			defer db.Close()
			if err != nil {
				slog.Error("error with seeding", slog.String("err", err.Error()))
			}
		}
		// -- set the db to use
		if exists.FileOrDir(sl.DB) {
			databases[sl.Table] = sl.DB
		}

	}
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
