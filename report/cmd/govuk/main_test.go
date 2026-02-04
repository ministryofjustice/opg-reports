package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domain/govuk/govuk"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"
	"testing"

	"github.com/google/go-github/v81/github"
)

func TestGovUKWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
	)

	if os.Getenv("GH_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GH_TOKEN"))

		err = downloadAssets(ctx, log, client.Repositories, &govuk.Options{
			Tag:       govUKReleaseTag,
			Directory: dir,
		})

		if err != nil {
			t.Errorf("unexpected import error: [%s]", err.Error())
			t.FailNow()
		}
	} else {
		t.SkipNow()
	}

}
