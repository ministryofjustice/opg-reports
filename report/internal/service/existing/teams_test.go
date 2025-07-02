package existing

import (
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/githubr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestTeamsInsert(t *testing.T) {
	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./existing-teams.db")
	conf.Github.Metadata.Asset = "test_accounts_v1.json"

	gh, _ := githubr.New(ctx, log, conf)
	sq, _ := sqlr.New(ctx, log, conf)
	// existing srv
	srv, err := New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	stmts, err := srv.InsertTeams(&mockedGitHubClient{}, gh, sq)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(stmts) <= 0 {
		t.Errorf("inserts failed")
	}

}

func TestTeamsgetTeamsFromMetadata(t *testing.T) {
	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	srv, _ := New(ctx, log, conf)
	gh, _ := githubr.New(ctx, log, conf)

	teams, err := srv.getTeamsFromMetadata(&mockedGitHubClient{}, gh, &teamDownloadOptions{
		Owner:      "testowner",
		Repository: "test-repo",
		AssetName:  "test_accounts_v1.json",
		UseRegex:   false,
		Dir:        dir,
	})

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(teams) != 1 {
		t.Error("unexpected team data returned")
		utils.Debug(teams)
	}

}
