package awscost_test

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/awscost"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestServiceNew(t *testing.T) {
	var (
		err error
		srv *awscost.Service[*awscost.AwsCost]
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awscost-connection.db")

	lg := utils.Logger("ERROR", "TEXT")
	rep, _ := sqlr.NewWithSelect[*awscost.AwsCost](ctx, lg, cfg)

	srv, err = awscost.NewService(ctx, lg, cfg, rep)
	if err != nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
	}
	defer srv.Close()

	srv, err = awscost.NewService[*awscost.AwsCost](ctx, nil, nil, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a log or repository")
	}
	defer srv.Close()

	srv, err = awscost.NewService[*awscost.AwsCost](ctx, lg, nil, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a repository")
	}
	defer srv.Close()

}

func TestServiceGetAll(t *testing.T) {
	var (
		err  error
		dir  = t.TempDir()
		ctx  = t.Context()
		conf = config.NewConfig()
		log  = utils.Logger("ERROR", "TEXT")
	)
	conf.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awscost-getall.db")

	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	seeder.Teams(sqc)
	seeder.AwsAccounts(sqc)
	inserts, _ := seeder.AwsCosts(sqc)

	// generate the service useing default
	srv := awscost.Default[*awscost.AwsCost](ctx, log, conf)
	if srv == nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
		t.FailNow()
	}
	defer srv.Close()

	all, err := srv.GetAll()
	if err != nil {
		t.Errorf("unexpected error getting data: [%s]", err.Error())
	}

	if len(all) != len(inserts) {
		t.Errorf("mismatched number of records returned compared to insert")
	}

}
