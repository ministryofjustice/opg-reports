package existing

import (
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/githubr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type tTeam struct{}

func TestTeamsInsert(t *testing.T) {
	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	conf.Database.Path = filepath.Join(dir, "./existing-teams.db")

	gh, _ := githubr.New(ctx, log, conf)
	sq, _ := sqlr.New(ctx, log, conf)
	// existing srv
	srv, err := New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	stmts, err := srv.InsertTeams(gh, sq)
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

	_, err = srv.getTeamsFromMetadata(gh, &teamDownloadOptions{
		Owner:      "ministryofjustice",
		Repository: "opg-metadata",
		AssetName:  "metadata.tar.gz",
		UseRegex:   false,
		Dir:        dir,
	})

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

}
