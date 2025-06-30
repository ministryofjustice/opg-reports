package metadata

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type testM struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestOpgMetaDataServiceDownloadAndExtract(t *testing.T) {
	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	if conf.Github.Token == "" {
		t.Skip("No GITHUB_TOKEN, skipping test")
	}

	ghs, _ := gh.New(ctx, log, conf)
	srv, _ := NewService[*testM](ctx, log, conf, ghs)
	if srv == nil {
		t.Errorf("error making service")
		t.FailNow()
	}
	defer srv.Close()

	srv.SetDirectory(dir)
	local, err := srv.DownloadAndExtractAsset("gitleaks", "gitleaks", "gitleaks_(.*)_darwin_arm64.tar.gz", true)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !utils.DirExists(local) {
		t.Errorf("asset was not extracted")
		t.FailNow()
	}

	files := utils.FileList(local, "")
	if len(files) <= 0 {
		t.Errorf("no files were extracted")
	}

}

// TestOpgMetaDataServiceDownloadAndReturn tests that json file is read from the
// downloaded releases asset correctly, so this still uses the known opg repo
// to do that
// That means this requires an active github token to access the private repo
func TestOpgMetaDataServiceDownloadAndReturn(t *testing.T) {
	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	if conf.Github.Token == "" {
		t.Skip("No GITHUB_TOKEN, skipping test")
	}

	ghs, _ := gh.New(ctx, log, conf)
	srv, _ := NewService[*testM](ctx, log, conf, ghs)
	if srv == nil {
		t.Errorf("error making service")
		t.FailNow()
	}
	defer srv.Close()

	srv.SetDirectory(dir)
	data, err := srv.DownloadAndReturn("ministryofjustice", "opg-metadata", "metadata.tar.gz", false, "accounts.json")

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(data) <= 0 {
		t.Errorf("expected data to be found")
	}

}
