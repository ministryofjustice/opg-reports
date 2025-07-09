package seed

import (
	"path/filepath"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

func TestSeedServiceTeam(t *testing.T) {

	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./seed-teams.db")

	sq := sqlr.Default(ctx, log, conf)
	seeder := Default(ctx, log, conf)

	res, err := seeder.Teams(sq)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(res) != len(teamSeeds) {
		t.Errorf("number of inserted records doesnt match number of seeds")
	}
}
