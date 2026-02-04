package upload

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domain/downloads/download"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/logger"
	"opg-reports/report/internal/utils/times"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestDomainUploadWithoutMock(t *testing.T) {

	var (
		err       error
		client    *s3.Client
		dlFile    string
		dlContent []byte
		data      string          = times.AsString(time.Now(), times.FULL)
		ctx       context.Context = t.Context()
		log       *slog.Logger    = logger.New("error")
		upDir     string          = t.TempDir()
		downDir   string          = t.TempDir()

		opts *Options = &Options{
			Bucket:   "report-data-development",
			Key:      "test/test.txt",
			Filepath: filepath.Join(upDir, "test.txt"),
		}
	)
	// skip if no aws values set
	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		t.SkipNow()
	}
	// write data to tmp file for testing
	// os.MkdirAll(filepath.Dir(opts.Filepath), os.ModePerm)
	err = os.WriteFile(opts.Filepath, []byte(data), os.ModePerm)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

	client, err = awsclients.New[*s3.Client](ctx, log, "eu-west-1")
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

	_, err = UploadItemToBucket(ctx, log, client, opts)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

	// now download the file to compare content
	dlFile, err = download.DownloadItemFromBucket(ctx, log, client, &download.Options{
		Bucket:    opts.Bucket,
		Key:       opts.Key,
		Directory: downDir,
	})
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}

	dlContent, err = os.ReadFile(dlFile)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
		t.FailNow()
	}
	if string(dlContent) != data {
		t.Errorf("downloaded content does not match uploaded.")
	}

}
