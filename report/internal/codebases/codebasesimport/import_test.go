package codebasesimport

import (
	"context"
	"opg-reports/report/internal/global/migrations"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/ghclients"
	"opg-reports/report/package/logger"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v81/github"
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
	err = Import(ctx, clients, &Args{
		DB:         dbpath,
		Driver:     "sqlite3",
		OrgSlug:    "ministryofjustice",
		ParentSlug: "opg",
	})
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

}
