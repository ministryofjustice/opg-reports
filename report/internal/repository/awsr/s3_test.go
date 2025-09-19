package awsr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockRepositoryS3BucketLister struct{}

func (self *mockRepositoryS3BucketLister) ListBucket(client s3.ListObjectsV2APIClient, bucket string, prefix string) ([]string, error) {
	return []string{
		fmt.Sprintf("%s/%s%s", bucket, prefix, "sample-00.json"),
		fmt.Sprintf("%s/%s%s", bucket, prefix, "sample-01.json"),
		fmt.Sprintf("%s/%s%s", bucket, prefix, "sample-01.csv"),
	}, nil
}

// mockedRepositoryS3BucketDownloader provides a mocked version of DownloadBucket that writes a dummy cost file to a
// known location and returns that as the file path
type mockedRepositoryS3BucketDownloader struct{}

// DownloadBucket generates a file with dummy cost data in to for testing inserts
func (self *mockedRepositoryS3BucketDownloader) DownloadBucket(client ClientS3ListAndGetter, bucket string, prefix string, directory string) (downloaded []string, err error) {
	var file = filepath.Join(directory, "sample-costs.json")
	var content = `[
	{
		"id": 0,
		"ts": "2024-08-15 18:52:55.055478 +0000 UTC",
		"organisation": "OPG",
		"account_id": "001A",
		"account_name": "Account 1A",
		"unit": "TEAM-A",
		"label": "A",
		"environment": "development",
		"service": "Amazon Simple Storage Service",
		"region": "eu-west-1",
		"date": "2025-05-31",
		"cost": "0.2309542206"
	},
	{
		"id": 0,
		"ts": "2024-08-15 18:52:55.055478 +0000 UTC",
		"organisation": "OPG",
		"account_id": "001A",
		"account_name": "Account 1A",
		"unit": "TEAM-A",
		"label": "A",
		"environment": "development",
		"service": "Amazon Simple Storage Service",
		"region": "eu-west-1",
		"date": "2025-04-31",
		"cost": "107.53"
	}
]`
	err = os.WriteFile(file, []byte(content), os.ModePerm)
	downloaded = append(downloaded, file)
	return
}

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
