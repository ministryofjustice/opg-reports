package awsr

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestS3BucketList(t *testing.T) {
	var err error

	repository := &mockRepositoryS3BucketLister{}

	files, err := repository.ListBucket(nil, "test-bucket-name", "prefix/")
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
	)
	os.MkdirAll(downloadTo, os.ModePerm)

	repository := &mockedRepositoryS3BucketDownloader{}

	files, err := repository.DownloadBucket(nil, "test-bucket-name", "prefix/", downloadTo)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	if len(files) <= 0 {
		t.Errorf("failed to list items in bucket")
	}
}

func TestS3BucketUpload(t *testing.T) {
	var (
		err        error
		ctx        = context.TODO()
		conf       = config.NewConfig()
		log        = utils.Logger("WARN", "TEXT")
		dir        = t.TempDir()
		sampleFile = filepath.Join(dir, "test.json")
	)
	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		t.SkipNow()
	} else {

		os.WriteFile(sampleFile, []byte(`test file!`), os.ModePerm)
		repo := Default(ctx, log, conf)
		client := DefaultClient[*s3.Client](ctx, "eu-west-1")
		_, err = repo.UploadItemToBucket(client, "report-data-development", "test/test.json", sampleFile)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}

}
