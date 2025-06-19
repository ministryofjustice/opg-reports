package opgmetadata

import (
	"log/slog"
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestOpgMetaDataServiceDownload(t *testing.T) {

	if utils.GetEnvVar("GH_TOKEN", "") == "" {
		t.Skip("No GH_TOKEN, skipping test")
	}

	var (
		err error
		dir string = t.TempDir()
		ctx        = t.Context()
		cfg        = config.NewConfig()
		lg         = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)

	cfg.Github.Organisation = "ministryofjustice"
	gh, _ := gh.New(ctx, lg, cfg)
	srv, _ := NewService(ctx, lg, cfg, gh)

	srv.SetDirectory(dir)
	err = srv.Download()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

}

func TestOpgMetaDataServiceGetAccounts(t *testing.T) {

	if utils.GetEnvVar("GH_TOKEN", "") == "" {
		t.Skip("No GH_TOKEN, skipping test")
	}

	var (
		err error
		dir string = t.TempDir()
		ctx        = t.Context()
		cfg        = config.NewConfig()
		lg         = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)

	cfg.Github.Organisation = "ministryofjustice"
	gh, _ := gh.New(ctx, lg, cfg)
	srv, _ := NewService(ctx, lg, cfg, gh)
	srv.SetDirectory(dir)

	accs, err := srv.GetAllAccounts()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(accs) <= 0 {
		t.Errorf("error with number of accounts found")
	}

}

func TestOpgMetaDataServiceGetTeams(t *testing.T) {

	if utils.GetEnvVar("GH_TOKEN", "") == "" {
		t.Skip("No GH_TOKEN, skipping test")
	}

	var (
		err error
		dir string = t.TempDir()
		ctx        = t.Context()
		cfg        = config.NewConfig()
		lg         = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)

	cfg.Github.Organisation = "ministryofjustice"
	gh, _ := gh.New(ctx, lg, cfg)
	srv, _ := NewService(ctx, lg, cfg, gh)
	srv.SetDirectory(dir)

	teams, err := srv.GetAllTeams()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(teams) <= 0 {
		t.Errorf("error with number of teams found")
	}

}
