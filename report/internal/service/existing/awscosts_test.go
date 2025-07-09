package existing

import (
	"path/filepath"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/seed"
	"opg-reports/report/internal/utils"
)

func TestAwsCostsInsert(t *testing.T) {
	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("INFO", "TEXT")
	)

	// set config values
	conf.Database.Path = filepath.Join(dir, "./existing-awscosts.db")

	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	seeder.Teams(sqc)
	seeder.AwsAccounts(sqc)

	sq, _ := sqlr.New(ctx, log, conf)
	// existing srv
	srv, err := New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	stmts, err := srv.InsertAwsCosts(nil, &mockedRepositoryS3BucketDownloader{}, sq)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(stmts) <= 0 {
		t.Errorf("inserts failed")
	}

}
