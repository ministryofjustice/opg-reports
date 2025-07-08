package awsr

import (
	"os"
	"path/filepath"
	"testing"
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
