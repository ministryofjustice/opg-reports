package team

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/teams/teammodels"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"
	"testing"

	"github.com/google/go-github/v81/github"
)

func TestRedoTeamsWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		teams  []*teammodels.Team
	)

	if os.Getenv("GITHUB_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GITHUB_TOKEN"))
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		opts := &TeamDataOptions{
			Tag:           "v0.1.26",
			DataDirectory: dir,
		}
		teams, err = GetTeamData[*github.RepositoriesService](ctx, log, client.Repositories, opts)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		if len(teams) < 1 {
			t.Errorf("expected more teams in the list")
		}
	} else {
		t.SkipNow()
	}
}
