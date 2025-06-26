package s3bucket

import (
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

func TestS3BucketDownloadBucket(t *testing.T) {
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

	files, err := repo.ListBucket("report-data-development", "aws_costs/")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(files) <= 0 {
		t.Errorf("did not find any files in bucket")
	}

	dl, err := repo.Download("report-data-development", files, "__download/")
	if len(dl) <= 0 {
		t.Errorf("no files downloaded")
	}
	t.Fail()

}
