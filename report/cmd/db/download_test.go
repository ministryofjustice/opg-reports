package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domain/downloads/download"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/files"
	"opg-reports/report/internal/utils/logger"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// TestDBDownloadWithoutMock
// run with aws-vault exec shared-development-operator -- make test name="TestDBDownloadWithoutMock"
func TestDBDownloadWithoutMock(t *testing.T) {

	var (
		err    error
		client *s3.Client
		ctx    context.Context   = t.Context()
		log    *slog.Logger      = logger.New("error")
		dir    string            = t.TempDir()
		opts   *download.Options = &download.Options{
			Bucket:    "opg-reports-development",
			Key:       "database/api.db",
			Directory: dir,
		}
	)
	if os.Getenv("AWS_SESSION_TOKEN") == "" {
		t.SkipNow()
	}

	client, err = awsclients.New[*s3.Client](ctx, log, "eu-west-1")
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}

	err = downloadItem(ctx, log, client, opts)
	if err != nil {
		t.Errorf("unexpected error:\n%s", err.Error())
	}
	if !files.Exists(filepath.Join(dir, opts.Key)) {
		t.Errorf("file was not downloaded")
	}

}
