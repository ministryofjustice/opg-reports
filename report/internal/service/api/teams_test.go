package api

import (
	"path/filepath"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

func TestApiServiceGetAllTeams(t *testing.T) {

	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./api-get-all-teams.db")
	inserted, _ := seedDB(ctx, log, conf)

	store := sqlr.DefaultWithSelect[*Team](ctx, log, conf)
	service := Default[*Team](ctx, log, conf)

	teams, err := service.GetAllTeams(store)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(inserted.Teams) != len(teams) {
		t.Errorf("mismatching number of records: expected [%d] actual [%v]", len(inserted.Teams), len(teams))
	}

}
