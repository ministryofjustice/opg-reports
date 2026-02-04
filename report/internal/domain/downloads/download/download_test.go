package download

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/debugger"
	"opg-reports/report/internal/utils/logger"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestDomainDownloadWithoutMock(t *testing.T) {

	var (
		err    error
		dest   string
		client *s3.Client
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = "./s3" //t.TempDir()
		opts   *Options        = &Options{
			Bucket:    "opg-reports-development",
			Key:       "database/api.db",
			Directory: dir,
		}
	)

	if os.Getenv("AWS_SESSION_TOKEN") != "" {
		client, err = awsclients.New[*s3.Client](ctx, log, "eu-west-1")
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}

		dest, err = DownloadItemFromBucket(ctx, log, client, opts)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}

		debugger.Dump(dest)

	} else {
		t.SkipNow()
	}
	t.FailNow()
}
