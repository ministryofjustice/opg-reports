package codeowner

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/domain/codebases/codebasemodels"
	"opg-reports/report/internal/domain/codeowners/codeownermodels"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"
	"testing"

	"github.com/google/go-github/v81/github"
)

func TestDomainCodeownerWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		opts   *Input
		ctx    context.Context              = t.Context()
		log    *slog.Logger                 = logger.New("error")
		result []*codeownermodels.Codeowner = []*codeownermodels.Codeowner{}
	)
	opts = &Input{
		Codebases: []*codebasemodels.Codebase{
			{FullName: "ministryofjustice/opg-lpa", Name: "opg-lpa"},
			{FullName: "ministryofjustice/opg-use-an-lpa", Name: "opg-use-an-lpa"},
			{FullName: "ministryofjustice/opg-data-lpa-store", Name: "opg-data-lpa-store"},
		},
	}

	if os.Getenv("GITHUB_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GITHUB_TOKEN"))
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}

		result, err = GetCodeowners(ctx, log, client.Repositories, opts)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		if len(result) <= 0 {
			t.Errorf("unexpected count of results")
		}

	} else {
		t.SkipNow()
	}
}
