package existing

import (
	"path/filepath"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

func TestTeamsInsert(t *testing.T) {
	var (
		err    error
		dir    string = t.TempDir()
		ctx           = t.Context()
		conf          = config.NewConfig()
		log           = utils.Logger("ERROR", "TEXT")
		client        = &mockClientRepositoryReleaseListReleases{}
		// client = githubr.DefaultClient(conf).Repositories
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./existing-teams.db")
	conf.Metadata.ReleaseTag = "v1(.*)"
	conf.Metadata.AssetName = "asset.json"
	conf.Metadata.UseRegex = true

	gh, _ := githubr.New(ctx, log, conf)
	sq, _ := sqlr.New(ctx, log, conf)
	// existing srv
	srv, err := New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	stmts, err := srv.InsertTeams(client, gh, sq)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(stmts) <= 0 {
		t.Errorf("inserts failed")
	}

}
