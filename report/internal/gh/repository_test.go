package gh

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// TestGhAllReleases makes a call to a real api to check there are releases returned
// - will skip if no GH_TOKEN is set
func TestGhAllReleases(t *testing.T) {

	if utils.GetEnvVar("GH_TOKEN", "") == "" {
		t.Skip("No GH_TOKEN, skipping test")
	}

	var (
		err error
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)

	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	found, err := repo.GetReleases("ministryofjustice", "opg-reports", nil)
	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if len(found) <= 0 {
		t.Errorf("expected multiple releases to be returned")
	}

}

// TestGhAllReleases makes a call to a real api to check there are releases returned
// - will skip if no GH_TOKEN is set
func TestGhLastReleases(t *testing.T) {

	if utils.GetEnvVar("GH_TOKEN", "") == "" {
		t.Skip("No GH_TOKEN, skipping test")
	}

	var (
		err error
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)

	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	found, err := repo.GetLatestRelease("ministryofjustice", "opg-reports")
	if err != nil {
		t.Errorf("unexpected error found: %s", err.Error())
	}
	if found == nil {
		t.Errorf("no releases found")
	}
	fmt.Printf("%s\n", utils.MarshalStr(found))
	t.Fail()
}
