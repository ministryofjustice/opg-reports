package opgmetadata

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestOpgMetaDataServiceDownload(t *testing.T) {

	var (
		err error
		dir string = t.TempDir()
		ctx        = t.Context()
		cfg        = config.NewConfig()
		lg         = utils.Logger("WARN", "TEXT")
	)
	if cfg.Github.Token == "" {
		t.Skip("No GITHUB_TOKEN, skipping test")
	}

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

	var (
		err error
		dir string = t.TempDir()
		ctx        = t.Context()
		cfg        = config.NewConfig()
		lg         = utils.Logger("WARN", "TEXT")
	)
	if cfg.Github.Token == "" {
		t.Skip("No GITHUB_TOKEN, skipping test")
	}

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

	var (
		err error
		dir string = t.TempDir()
		ctx        = t.Context()
		cfg        = config.NewConfig()
		lg         = utils.Logger("WARN", "TEXT")
	)
	if cfg.Github.Token == "" {
		t.Skip("No GITHUB_TOKEN, skipping test")
	}

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
