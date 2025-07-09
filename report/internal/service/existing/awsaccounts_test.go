package existing

import (
	"path/filepath"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/seed"
	"opg-reports/report/internal/utils"
)

func TestAwsAccountsInsert(t *testing.T) {
	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./existing-awsaccounts.db")
	conf.Metadata.Asset = "test_accounts_v1.json"

	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	seeder.Teams(sqc)

	gh, _ := githubr.New(ctx, log, conf)
	sq, _ := sqlr.New(ctx, log, conf)
	// existing srv
	srv, err := New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	stmts, err := srv.InsertAwsAccounts(&mockedGitHubClient{}, gh, sq)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(stmts) <= 0 {
		t.Errorf("inserts failed")
	}

}
