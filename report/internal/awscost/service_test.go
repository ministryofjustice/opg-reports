package awscost_test

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/awscost"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
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

	lg := utils.Logger("WARN", "TEXT")
	rep, _ := sqldb.New[*awscost.AwsCost](ctx, lg, cfg)

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
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("WARN", "TEXT")
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awscost-getall.db")

	// seed both teams and accounts before costs
	team.Seed(ctx, lg, cfg, nil)
	awsaccount.Seed(ctx, lg, cfg, nil)

	inserts, err := awscost.Seed(ctx, lg, cfg, nil)
	if err != nil {
		t.Errorf("unexpected error seeding: [%s]", err.Error())
	}

	// generate the service useing default
	srv, err := awscost.Default[*awscost.AwsCost](ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
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
