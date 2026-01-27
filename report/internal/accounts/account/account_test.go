package account

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/accounts/accountmodels"
	"opg-reports/report/internal/utils/ghclients"
	"opg-reports/report/internal/utils/logger"
	"os"
	"testing"

	"github.com/google/go-github/v81/github"
)

func TestRedoAccountsWithoutMock(t *testing.T) {

	var (
		err    error
		client *github.Client
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		dir    string          = t.TempDir()
		data   []*accountmodels.AwsAccountImport
	)

	if os.Getenv("GITHUB_TOKEN") != "" {
		client, err = ghclients.New(ctx, log, os.Getenv("GITHUB_TOKEN"))
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		opts := &GetAwsAccountDataOptions{
			Tag:           "v0.1.26",
			DataDirectory: dir,
		}

		data, err = GetAwsAccountData[*github.RepositoriesService](ctx, log, client.Repositories, opts)
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		if len(data) < 30 {
			t.Errorf("expected more teams in the list")
		}
	} else {
		t.SkipNow()
	}
}
