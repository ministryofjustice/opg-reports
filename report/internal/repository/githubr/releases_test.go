package githubr

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// mockedR is a dummy stuck used for testing functions by using fixed values
// so we dont have to makle api calls to github during testing
type mockedR struct{}

var t = true
var f = false
var id int64 = 2003
var ids = []int64{1000, 1001, 1002, 1003}
var nm = "test_v1_darwin_arm64.txt"
var nms = []string{"R00", "R01", "R02", "R03"}
var rel = &github.RepositoryRelease{
	ID:         &ids[3],
	Name:       &nms[3],
	Draft:      &f,
	Prerelease: &t,
	Assets: []*github.ReleaseAsset{
		{ID: &id, Name: &nm},
	},
}
var rels = []*github.RepositoryRelease{
	{ID: &ids[0], Name: &nms[0], Draft: &t, Prerelease: &t, Assets: []*github.ReleaseAsset{}},
	{ID: &ids[1], Name: &nms[1], Draft: &f, Prerelease: &f, Assets: []*github.ReleaseAsset{}},
	{ID: &ids[2], Name: &nms[2], Draft: &f, Prerelease: &t, Assets: []*github.ReleaseAsset{}},
	rel,
}

func (self *mockedR) ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) (all []*github.RepositoryRelease, resp *github.Response, err error) {
	all = rels
	resp = &github.Response{
		NextPage: 0,
	}
	return
}

func (self *mockedR) GetLatestRelease(ctx context.Context, owner, repo string) (release *github.RepositoryRelease, resp *github.Response, err error) {
	release = rel
	return
}

// DownloadReleaseAsset is a dummy function that mimics downloading a specific id to a known file localtion
func (self *mockedR) DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error) {
	var fp = filepath.Join("./", nm)
	err = os.WriteFile(fp, []byte(`test-file`), os.ModePerm)
	rc, _ = os.Open(fp)
	return
}

// TestGithubrAllReleases makes a call to a real api to check there are releases returned
// Uses a public repo to reduce need for token auth
//
// Warning - makes a real api call to the github api
func TestGithubrAllReleases(t *testing.T) {

	var (
		err error
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)
	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	found, err := repo.GetReleases(&mockedR{}, "", "", nil)
	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if len(found) <= 0 {
		t.Errorf("expected multiple releases to be returned")
	}
	// check filtering on the release set
	nopre, err := repo.GetReleases(&mockedR{}, "", "", &GetReleaseOptions{ExcludePrereleases: true})
	ok := true
	for _, r := range nopre {
		if *r.Prerelease == true {
			ok = false
		}
	}
	if !ok {
		t.Errorf("got an unexpected pre-release")
	}

	nodr, err := repo.GetReleases(&mockedR{}, "", "", &GetReleaseOptions{ExcludeDraft: true})
	ok = true
	for _, r := range nodr {
		if *r.Draft == true {
			ok = false
		}
	}
	if !ok {
		t.Errorf("got an unexpected draft")
	}
	noa, err := repo.GetReleases(&mockedR{}, "", "", &GetReleaseOptions{ExcludeNoAssets: true})
	ok = true
	for _, r := range noa {
		if r.Assets == nil {
			ok = false
		}
	}
	if !ok {
		t.Errorf("got an unexpected empty asset list")
	}

}

// TestGithubrLatestRelease
func TestGithubrLatestRelease(t *testing.T) {

	var (
		err error
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)

	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	found, err := repo.GetLatestRelease(&mockedR{}, "", "")
	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if found == nil {
		t.Errorf("no releases found")
	}
}

// TestGithubrLatestReleaseAssetAndDownload
func TestGithubrLatestReleaseAssetAndDownload(t *testing.T) {

	var (
		err  error
		dir  = t.TempDir()
		dest = fmt.Sprintf("%s/%s", dir, "test.txt")
		ctx  = t.Context()
		cfg  = config.NewConfig()
		lg   = utils.Logger("ERROR", "TEXT")
	)
	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	found, err := repo.GetLatestReleaseAsset(
		&mockedR{},
		"", "", "test_(.*)_darwin_arm64.txt",
		true)

	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
		t.FailNow()
	}
	if found == nil {
		t.Errorf("no releases found")
		t.FailNow()
	}
	// this is the test calue that should come back
	if *found.Name != "test_v1_darwin_arm64.txt" {
		t.Errorf("no releases match found")
	}

	f, err := repo.DownloadReleaseAsset(
		&mockedR{},
		"gitleaks", "gitleaks", *found.ID,
		dest)
	defer f.Close()

	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if f == nil {
		t.Errorf("file pointer is nil")
	}
	if !utils.FileExists(dest) {
		t.Errorf("file missing")
	}
	// remove the dummy file
	os.RemoveAll(filepath.Join("./", nm))

}

// TestGithubrLastReleaseAssetByName
func TestGithubrLastReleaseAssetByName(t *testing.T) {

	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	path, err := repo.DownloadReleaseAssetByName(
		&mockedR{},
		"", "", "test_(.*)_darwin_arm64.txt",
		true, dir)

	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}

	if path == "" {
		t.Errorf("file path is nil")
	}
	os.RemoveAll(filepath.Join("./", nm))

}
