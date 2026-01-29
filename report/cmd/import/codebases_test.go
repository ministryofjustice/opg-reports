package main

import (
	"context"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/domain/codebases/codebaseselects"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/ptr"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
)

// mockCodebaseClient is a mocked client that returns successful results
type mockCodebaseClient struct{}

func (self *mockCodebaseClient) ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) (repos []*github.Repository, resp *github.Response, err error) {

	repos = []*github.Repository{
		{
			Name:     ptr.Ptr("mock-repo-a"),
			FullName: ptr.Ptr("ministryofjustice/mock-repo-a"),
			HTMLURL:  ptr.Ptr("https://test.local/ministryofjustice/mock-repo-a"),
			Archived: ptr.Ptr(false),
		},
		{
			Name:     ptr.Ptr("mock-repo-b"),
			FullName: ptr.Ptr("ministryofjustice/mock-repo-b"),
			HTMLURL:  ptr.Ptr("https://test.local/ministryofjustice/mock-repo-b"),
			Archived: ptr.Ptr(false),
		},
		{
			Name:     ptr.Ptr("mock-repo-c"),
			FullName: ptr.Ptr("ministryofjustice/mock-repo-c"),
			HTMLURL:  ptr.Ptr("https://test.local/ministryofjustice/mock-repo-c"),
			Archived: ptr.Ptr(true),
		},
	}
	resp = &github.Response{
		NextPage: 0,
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
	}
	return
}

func TestImportsCodebasesWithMock(t *testing.T) {

	var (
		err    error
		db     *sqlx.DB
		client *mockCodebaseClient = &mockCodebaseClient{}
		ctx    context.Context     = t.Context()
		log    *slog.Logger        = logger.New("error")
		dir    string              = t.TempDir()
		dbPath string              = filepath.Join(dir, "test-import-mock-codebases.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbmigrations.Migrate(ctx, log, db)
	defer db.Close()

	err = importCodebases(ctx, log, client, db)
	if err != nil {
		t.Errorf("unexpected import error: [%s]", err.Error())
		t.FailNow()
	}

	data, err := codebaseselects.All(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected select error: [%s]", err.Error())
		t.FailNow()
	}
	if len(data) != 2 {
		t.Errorf("expected exactly 2 repos from the mock data to be created.")
	}

}

func TestImportsCodebasesWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		db     *sqlx.DB
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-codebases.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbmigrations.Migrate(ctx, log, db)
	defer db.Close()

	if os.Getenv("GH_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GH_TOKEN"))

		err = importCodebases(ctx, log, client.Teams, db)
		if err != nil {
			t.Errorf("unexpected import error: [%s]", err.Error())
			t.FailNow()
		}
	} else {
		t.SkipNow()
	}
}
