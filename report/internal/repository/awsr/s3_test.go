package awsr

import (
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type tS3 struct{}

func TestS3BucketList(t *testing.T) {
	var (
		err        error
		repository *Repository
		ctx        = t.Context()
		conf       = config.NewConfig()
		log        = utils.Logger("ERROR", "TEXT")
	)

	if conf.Aws.GetToken() == "" {
		t.Skip("No AWS_SESSION_TOKEN, skipping test")
	}

	repository, err = New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	files, err := repository.ListBucket("report-data-development", "github_standards/")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	if len(files) <= 0 {
		t.Errorf("failed to list items in bucket")
	}

}

func TestS3BucketDownload(t *testing.T) {
	var (
		err        error
		dir        = t.TempDir()
		downloadTo = filepath.Join(dir, "__download/")
		repository *Repository
		ctx        = t.Context()
		conf       = config.NewConfig()
		log        = utils.Logger("ERROR", "TEXT")
	)

	if conf.Aws.GetToken() == "" {
		t.Skip("No AWS_SESSION_TOKEN, skipping test")
	}

	repository, err = New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	files, err := repository.DownloadBucket("report-data-development", "github_standards/", downloadTo)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	if len(files) <= 0 {
		t.Errorf("failed to list items in bucket")
	}
}
