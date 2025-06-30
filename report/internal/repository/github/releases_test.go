package github

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// TestGhAllReleases makes a call to a real api to check there are releases returned
// Uses a public repo to reduce need for token auth
//
// Warning - makes a real api call to the github api
func TestGhAllReleases(t *testing.T) {

	var (
		err error
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)
	if cfg.Github.Token == "" {
		t.Skip("No GITHUB_TOKEN, skipping test")
	}

	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	found, err := repo.GetReleases("actions", "checkout", nil)
	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if len(found) <= 0 {
		t.Errorf("expected multiple releases to be returned")
	}

}

// TestGhLastReleases
// Uses a public repo to reduce need for token auth
// Warning - makes a real api call to the github api
func TestGhLastReleases(t *testing.T) {

	var (
		err error
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)
	if cfg.Github.Token == "" {
		t.Skip("No GITHUB_TOKEN, skipping test")
	}

	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	found, err := repo.GetLatestRelease("actions", "checkout")
	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if found == nil {
		t.Errorf("no releases found")
	}
}

// TestGhLastReleaseAssetAndDownload
// Use a poplar repo that generates a tarball (gitleaks) to
// test a download
// Uses a public repo to reduce need for token auth
// Warning - makes a real api call to the github api
func TestGhLastReleaseAssetAndDownload(t *testing.T) {

	var (
		err error
		dir = t.TempDir()
		fp  = fmt.Sprintf("%s/%s", dir, "gitleaks.tar.gz")
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)
	if cfg.Github.Token == "" {
		t.Skip("No GITHUB_TOKEN, skipping test")
	}

	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	found, err := repo.GetLatestReleaseAsset("gitleaks", "gitleaks", "gitleaks_(.*)_darwin_arm64.tar.gz", true)
	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if found == nil {
		t.Errorf("no releases found")
		t.FailNow()
	}

	f, err := repo.DownloadReleaseAsset("gitleaks", "gitleaks", *found.ID, fp)
	defer f.Close()

	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if f == nil {
		t.Errorf("file pointer is nil")
	}
	if !utils.FileExists(fp) {
		t.Errorf("file missing")
	}

}
