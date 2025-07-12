package front

import (
	"context"
	"opg-reports/report/config"
	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/utils"
	"path/filepath"
	"testing"
)

func TestFrontServiceDownloadGovUKFrontEnd(t *testing.T) {
	var (
		err    error
		dir    = filepath.Join(t.TempDir(), "govuk")
		ctx    = context.TODO()
		log    = utils.Logger("ERROR", "TEXT")
		conf   = config.NewConfig()
		client = githubr.DefaultClient(conf).Repositories
		store  = githubr.Default(ctx, log, conf)
		serv   = Default(ctx, log, conf)
	)

	files, _, err := serv.DownloadGovUKFrontEnd(client, store, dir)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(files) <= 0 {
		t.Errorf("files not found in zip")
	}
}
