package govuk

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"
	"testing"

	"github.com/google/go-github/v81/github"
)

func TestDomainGovUKDLWithoutMock(t *testing.T) {

	var (
		err    error
		dest   string
		client *github.Client
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = "./tmp/" //t.TempDir()
	)

	if os.Getenv("GH_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GH_TOKEN"))
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		opts := &Options{
			Tag:       "v5.11.0",
			Directory: dir,
		}
		dest, err = DownloadFrontEnd(ctx, log, client.Repositories, opts)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		if dest != dir {
			t.Error("extracted zip into different location that set")
		}

	} else {
		t.SkipNow()
	}
}
