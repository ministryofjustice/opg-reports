package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codeowners/codeowner"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/ptr"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
)

type mockCodeownerClient struct{}

func (self *mockCodeownerClient) ListTeams(ctx context.Context, owner, repo string, opts *github.ListOptions) (teams []*github.Team, resp *github.Response, err error) {
	var requiredTeamParent = "opg"

	parent := &github.Team{
		Slug: ptr.Ptr(requiredTeamParent),
	}

	teams = []*github.Team{
		{
			Name:    ptr.Ptr("mock-other"),
			Slug:    ptr.Ptr("mock-other"),
			HTMLURL: ptr.Ptr("https://github.com/orgs/" + owner + "/mock-other"),
		},
		{
			Name:    ptr.Ptr("mock-team-a"),
			Slug:    ptr.Ptr("mock-team-a"),
			HTMLURL: ptr.Ptr("https://github.com/orgs/" + owner + "/mock-team-a"),
			Parent:  parent,
		},
		{
			Name:    ptr.Ptr("mock-team-b"),
			Slug:    ptr.Ptr("mock-team-b"),
			HTMLURL: ptr.Ptr("https://github.com/orgs/" + owner + "/mock-team-b"),
			Parent:  parent,
		},
	}
	// add extra team for mock-repo-b
	if repo == "mock-repo-b" {
		teams = append(teams, &github.Team{
			Name:    ptr.Ptr("mock-team-c"),
			Slug:    ptr.Ptr("mock-team-c"),
			HTMLURL: ptr.Ptr("https://github.com/orgs/" + owner + "/mock-team-c"),
			Parent:  parent,
		})
	}

	resp = &github.Response{
		NextPage: 0,
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
	}
	return
}

func (self *mockCodeownerClient) DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (buff io.ReadCloser, resp *github.Response, err error) {

	content := fmt.Sprintf(`* @%s/mock-team-a @%s/mock-team-b%s`, owner, owner, "\n")
	if repo == "mock-repo-b" {
		content += fmt.Sprintf(`.github/  @%s/mock-team-codeowner%s`, owner, "\n")
	}
	buff = io.NopCloser(bytes.NewBuffer([]byte(content)))

	resp = &github.Response{
		NextPage: 0,
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
	}
	return
}

func TestCMDImportsCodeownersWithMock(t *testing.T) {

	var (
		err    error
		db     *sqlx.DB
		client *mockCodeownerClient = &mockCodeownerClient{}
		code   []*codebasemodels.Codebase
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-mock-codeowners.db")
		org    string          = "mock-org"
	)
	code = []*codebasemodels.Codebase{
		{FullName: org + "/mock-repo-a", Name: "mock-repo-a"},
		{FullName: org + "/mock-repo-b", Name: "mock-repo-b"},
		{FullName: org + "/mock-repo-c", Name: "mock-repo-c"},
	}
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbsetup.Migrate(ctx, log, db)
	defer db.Close()

	err = importCodeowners(ctx, log, client, db, &codeowner.Input{
		Codebases:  code,
		ParentTeam: "opg",
		OrgSlug:    org,
	})

	if err != nil {
		t.Errorf("unexpected import error: [%s]", err.Error())
		t.FailNow()
	}

}

func TestCMDImportsCodeownersWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		db     *sqlx.DB
		code   []*codebasemodels.Codebase
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-codeowners.db")
	)

	if os.Getenv("GH_TOKEN") == "" {
		t.SkipNow()
	}

	code = []*codebasemodels.Codebase{
		{FullName: "ministryofjustice/opg-lpa", Name: "opg-lpa"},
		{FullName: "ministryofjustice/opg-use-an-lpa", Name: "opg-use-an-lpa"},
		{FullName: "ministryofjustice/opg-data-lpa-store", Name: "opg-data-lpa-store"},
	}
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbsetup.Migrate(ctx, log, db)
	defer db.Close()

	client, err = ghclients.New(ctx, log, os.Getenv("GH_TOKEN"))
	err = importCodeowners(ctx, log, client.Repositories, db, &codeowner.Input{
		Codebases:  code,
		ParentTeam: "opg",
		OrgSlug:    "ministryofjustice",
	})
	if err != nil {
		t.Errorf("unexpected import error: [%s]", err.Error())
		t.FailNow()
	}

}
