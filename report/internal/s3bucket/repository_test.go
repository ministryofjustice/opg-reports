package s3bucket

import (
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// TestS3BucketListBucket uses real config values to try and list files in a known bucket
func TestS3BucketListBucket(t *testing.T) {
	var (
		err error
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("WARN", "TEXT")
	)

	if cfg.Aws.GetToken() == "" {
		t.Skip("No AWS_SESSION_TOKEN, skipping test")
	}

	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	files, err := repo.ListBucket("report-data-development", "github_standards/")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(files) <= 0 {
		t.Errorf("did not find any files in bucket")
	}

}

// TestS3BucketDownloadBucket makes an active s3 call and downloads some files,
// making it slow
func TestS3BucketListAndDownloadBucket(t *testing.T) {
	var (
		err   error
		dir   = t.TempDir()
		dlDir = filepath.Join(dir, "__download/")
		ctx   = t.Context()
		cfg   = config.NewConfig()
		lg    = utils.Logger("WARN", "TEXT")
	)

	if cfg.Aws.GetToken() == "" {
		t.Skip("No AWS_SESSION_TOKEN, skipping test")
	}

	repo, err := New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	files, err := repo.ListBucket("report-data-development", "github_standards/")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(files) <= 0 {
		t.Errorf("did not find any files in bucket")
	}

	dl, err := repo.Download("report-data-development", files, dlDir)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(dl) != len(files) {
		t.Errorf("no files downloaded")
	}

}
