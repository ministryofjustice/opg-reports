package githubr

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

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

	found, err := repo.GetReleases(&mockedClient{}, "", "", nil)
	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if len(found) <= 0 {
		t.Errorf("expected multiple releases to be returned")
	}
	// check filtering on the release set
	nopre, err := repo.GetReleases(&mockedClient{}, "", "", &GetReleaseOptions{ExcludePrereleases: true})
	ok := true
	for _, r := range nopre {
		if *r.Prerelease == true {
			ok = false
		}
	}
	if !ok {
		t.Errorf("got an unexpected pre-release")
	}

	nodr, err := repo.GetReleases(&mockedClient{}, "", "", &GetReleaseOptions{ExcludeDraft: true})
	ok = true
	for _, r := range nodr {
		if *r.Draft == true {
			ok = false
		}
	}
	if !ok {
		t.Errorf("got an unexpected draft")
	}
	noa, err := repo.GetReleases(&mockedClient{}, "", "", &GetReleaseOptions{ExcludeNoAssets: true})
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

	found, err := repo.GetLatestRelease(&mockedClient{}, "", "")
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
		&mockedClient{},
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
		&mockedClient{},
		"gitleaks", "gitleaks", found,
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

	asset, path, err := repo.DownloadReleaseAssetByName(
		&mockedClient{},
		"", "", "test_(.*)_darwin_arm64.txt",
		true, dir)

	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if asset == nil {
		t.Errorf("unexpected nil asset returned: %s", err.Error())
	}

	if path == "" {
		t.Errorf("file path is nil")
	}

}
