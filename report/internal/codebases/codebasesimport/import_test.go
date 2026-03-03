package codebasesimport

import (
	"context"
	"opg-reports/report/internal/codebases/codebasesimport/args"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/ghclients"
	"opg-reports/report/package/logger"
	"opg-reports/report/package/times"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v84/github"
)

func TestCodebasesImportWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		ctx    context.Context = cntxt.AddLogger(t.Context(), logger.New("error"))
		dir    string          = t.TempDir()
		dbpath string          = filepath.Join(dir, "test-import.db")
	)
	if os.Getenv("GH_TOKEN") == "" {
		t.SkipNow()
	}

	client, err = ghclients.New(ctx, os.Getenv("GH_TOKEN"))
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

	migrations.Migrate(ctx, &migrations.Args{
		DB:     dbpath,
		Driver: "sqlite3",
	})

	clients := &Clients{
		Teams: client.Teams,
		Repos: client.Repositories,
	}

	err = Import(ctx, clients, &args.Args{
		DB:                dbpath,
		Driver:            "sqlite3",
		OrgSlug:           "ministryofjustice",
		ParentSlug:        "opg",
		IncludeStats:      false,
		IncludeCodeowners: false,
	})
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

}

func TestCodebasesImportStatsWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		ctx    context.Context = cntxt.AddLogger(t.Context(), logger.New("debug"))
		dir    string          = t.TempDir()
		dbpath string          = filepath.Join(dir, "test-import.db")
	)
	if os.Getenv("GH_TOKEN") == "" {
		t.SkipNow()
	}

	client, err = ghclients.New(ctx, os.Getenv("GH_TOKEN"))
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

	migrations.Migrate(ctx, &migrations.Args{
		DB:     dbpath,
		Driver: "sqlite3",
	})

	clients := &Clients{
		Teams: client.Teams,
		Repos: client.Repositories,
	}
	// just import a single repo for performance
	err = Import(ctx, clients, &args.Args{
		DB:           dbpath,
		Driver:       "sqlite3",
		OrgSlug:      "ministryofjustice",
		ParentSlug:   "opg",
		IncludeStats: true,
		FilterByName: "opg-lpa",
	})
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

}

func TestCodebasesImportMetricsWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		ctx    context.Context = cntxt.AddLogger(t.Context(), logger.New("debug"))
		dir    string          = t.TempDir()
		dbpath string          = filepath.Join(dir, "test-import.db")
	)
	if os.Getenv("GH_TOKEN") == "" {
		t.SkipNow()
	}

	client, err = ghclients.New(ctx, os.Getenv("GH_TOKEN"))
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

	migrations.Migrate(ctx, &migrations.Args{
		DB:     dbpath,
		Driver: "sqlite3",
	})

	clients := &Clients{
		Teams:   client.Teams,
		Repos:   client.Repositories,
		Actions: client.Actions,
	}
	// just import a single repo for performance
	err = Import(ctx, clients, &args.Args{
		DB:             dbpath,
		Driver:         "sqlite3",
		OrgSlug:        "ministryofjustice",
		ParentSlug:     "opg",
		IncludeMetrics: true,
		FilterByName:   "opg-use-an-lpa",
		DateStart:      times.MustFromString("2026-02-01"),
		DateEnd:        times.MustFromString("2026-03-03"),
	})
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
	t.FailNow()

}
