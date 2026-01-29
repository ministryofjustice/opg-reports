package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/domain/teams/teammodels"
	"opg-reports/report/internal/domain/teams/teamselects"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/marshal"
	"opg-reports/report/internal/utils/ptr"
	"opg-reports/report/internal/utils/zips"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v81/github"
	"github.com/jmoiron/sqlx"
)

// mockTeamsClient is a mocked client that returns successful results
type mockTeamsClient struct{}

// GetReleaseByTag mocked version to return preset version
func (self *mockTeamsClient) GetReleaseByTag(ctx context.Context, owner, repo, tag string) (rel *github.RepositoryRelease, resp *github.Response, err error) {
	var (
		mockTar = &github.ReleaseAsset{
			ID:          ptr.Ptr[int64](20001),
			Name:        ptr.Ptr("metadata.tar.gz"),
			ContentType: ptr.Ptr("application/tar+gzip"),
		}
		mockZip = &github.ReleaseAsset{
			ID:          ptr.Ptr[int64](20002),
			Name:        ptr.Ptr("metadata.zip"),
			ContentType: ptr.Ptr("application/zip"),
		}
	)
	resp = &github.Response{
		NextPage: 0,
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
	}
	rel = &github.RepositoryRelease{
		Name:       ptr.Ptr("mock-release"),
		Draft:      ptr.Ptr(false),
		Prerelease: ptr.Ptr(false),
		TagName:    ptr.Ptr(tag),
		Assets: []*github.ReleaseAsset{
			mockTar,
			mockZip,
		},
	}
	return
}

// DownloadReleaseAsset returns mocked version with the asset name as the content of the file
func (self *mockTeamsClient) DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error) {

	// create a temp zip file with teams.json file which contains dummy data
	// - write to tmp location, read and stream
	content := createDummyTeamZip()
	rc = io.NopCloser(bytes.NewBuffer(content))
	return
}

func TestImportsTeamsWithMock(t *testing.T) {

	var (
		err    error
		db     *sqlx.DB
		client *mockTeamsClient = &mockTeamsClient{}
		ctx    context.Context  = t.Context()
		log    *slog.Logger     = logger.New("error")
		dir    string           = t.TempDir()
		dbPath string           = filepath.Join(dir, "test-import-mock-teams.db")
	)
	// set the database
	db, err = dbconnection.Connection(ctx, log, "sqlite3", dbPath)
	if err != nil {
		t.Errorf("unexpected connection error: [%s]", err.Error())
		t.FailNow()
	}
	dbmigrations.Migrate(ctx, log, db)
	defer db.Close()

	err = teamsImport(ctx, log, client, db)
	if err != nil {
		t.Errorf("unexpected import error: [%s]", err.Error())
		t.FailNow()
	}

	data, err := teamselects.All(ctx, log, db)
	if err != nil {
		t.Errorf("unexpected select error: [%s]", err.Error())
		t.FailNow()
	}

	if len(data) <= 1 {
		t.Error("expected more results to be returned.")
	}

}

func TestImportsTeamsWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		db     *sqlx.DB
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		dbPath string          = filepath.Join(dir, "test-import-teams.db")
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

		err = teamsImport(ctx, log, client.Repositories, db)
		if err != nil {
			t.Errorf("unexpected import error: [%s]", err.Error())
			t.FailNow()
		}
	} else {
		t.SkipNow()
	}

}

// createDummyTeamZip makes fake accounts json file, then builds a zip
// and from that reads the zip and returns the content to simulate
// a downloaded file
func createDummyTeamZip() []byte {
	var mockAccounts = []*teammodels.Team{
		{Name: "mock-account-A"},
		{Name: "mock-account-B"},
	}

	tmpDir, _ := os.MkdirTemp("", "mock-team-*")
	// tmpDir := "./tmp"
	dataDir := filepath.Join(tmpDir, "src")
	os.MkdirAll(dataDir, os.ModePerm)
	// write the accounts to dummy file
	accountStr := marshal.ToString(mockAccounts)
	file := filepath.Join(dataDir, "teams.json")
	os.WriteFile(file, []byte(accountStr), os.ModePerm)

	// create the zip
	zipFile := filepath.Join(tmpDir, "metadata.zip")
	zips.Create(zipFile, []string{file}, dataDir)
	// read content and return
	content, _ := os.ReadFile(zipFile)
	// remove the tmp dir and its content
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	return content
}
