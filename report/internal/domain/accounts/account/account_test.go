package account

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"opg-reports/report/internal/domain/accounts/accountmodels"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/marshal"
	"opg-reports/report/internal/utils/ptr"
	"opg-reports/report/internal/utils/zips"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v81/github"
)

// mockGetter is a mocked client that returns successful results
type mockGetter struct{}

// GetReleaseByTag mocked version to return preset version
func (self *mockGetter) GetReleaseByTag(ctx context.Context, owner, repo, tag string) (rel *github.RepositoryRelease, resp *github.Response, err error) {
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
func (self *mockGetter) DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error) {

	// create a temp zip file with account.aws.json file which contains dummy data
	// - write to tmp location, read and stream
	content := createDummyZip()
	rc = io.NopCloser(bytes.NewBuffer(content))
	return
}

func TestDomainAccountsWithMock(t *testing.T) {

	var (
		err    error
		client *mockGetter     = &mockGetter{}
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		data   []*accountmodels.AwsAccount
	)

	opts := &GetAwsAccountDataOptions{
		Tag:           "v0.1.26",
		DataDirectory: dir,
	}

	data, err = GetAwsAccountData(ctx, log, client, opts)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
	if len(data) < 1 {
		t.Errorf("expected more teams in the list")
	}

}

func TestDomainAccountsWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		data   []*accountmodels.AwsAccount
	)

	if os.Getenv("GITHUB_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GITHUB_TOKEN"))
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		opts := &GetAwsAccountDataOptions{
			Tag:           "v0.1.26",
			DataDirectory: dir,
		}

		data, err = GetAwsAccountData[*github.RepositoriesService](ctx, log, client.Repositories, opts)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		if len(data) < 30 {
			t.Errorf("expected more teams in the list")
		}
	} else {
		t.SkipNow()
	}
}

// createDummyZip makes fake accounts json file, then builds a zip
// and from that reads the zip and returns the content to simulate
// a downloaded file
func createDummyZip() []byte {
	var mockAccounts = []*accountmodels.AwsAccount{
		{
			ID:          "mock-account-01",
			Name:        "mock-account",
			Label:       "mock-account",
			Environment: "development",
			TeamName:    "mock-team",
		},
	}

	tmpDir, _ := os.MkdirTemp("", "mock-account-*")
	// tmpDir := "./tmp"
	dataDir := filepath.Join(tmpDir, "src")
	os.MkdirAll(dataDir, os.ModePerm)
	// write the accounts to dummy file
	accountStr := marshal.ToString(mockAccounts)
	file := filepath.Join(dataDir, "accounts.aws.json")
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
