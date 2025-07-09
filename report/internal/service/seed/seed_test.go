package seed

import (
	"path/filepath"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

func TestSeedService(t *testing.T) {
	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./seed-test.db")
	sqc := sqlr.Default(ctx, log, conf)
	seeder := Default(ctx, log, conf)

	teams, err := seeder.Teams(sqc)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(teams) != len(teamSeeds) {
		t.Errorf("number of inserted records doesnt match number of seeds")
	}

	awsaccounts, err := seeder.AwsAccounts(sqc)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(awsaccounts) != len(awsAccountSeeds) {
		t.Errorf("number of inserted records doesnt match number of seeds")
	}

	awscosts, err := seeder.AwsCosts(sqc)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(awscosts) != len(awsCostSeeds) {
		t.Errorf("number of inserted records doesnt match number of seeds")
	}

}
