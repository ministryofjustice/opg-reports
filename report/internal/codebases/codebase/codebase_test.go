package codebase

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/codebases/codebasemodels"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"
	"testing"

	"github.com/google/go-github/v81/github"
)

func TestRedoCodebasesWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		data   []*codebasemodels.Codebase
	)

	if os.Getenv("GITHUB_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GITHUB_TOKEN"))
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		opts := &GetCodebasesOptions{ExcludeArchived: true}
		data, err = GetCodebases(ctx, log, client.Teams, opts)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		if len(data) < 5 {
			t.Errorf("expected more teams in the list")
		}

	} else {
		t.SkipNow()
	}
}
