package awsaccount_test

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/service/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/service/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestAwsAccountServiceNew(t *testing.T) {
	var (
		err error
		srv *awsaccount.Service[*awsaccount.AwsAccount]
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awsaccounts-connection.db")

	lg := utils.Logger("ERROR", "TEXT")
	rep, _ := sqldb.New[*awsaccount.AwsAccount](ctx, lg, cfg)

	srv, err = awsaccount.NewService(ctx, lg, cfg, rep)
	if err != nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
	}
	defer srv.Close()

	srv, err = awsaccount.NewService[*awsaccount.AwsAccount](ctx, nil, nil, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a log or repository")
	}
	defer srv.Close()

	srv, err = awsaccount.NewService[*awsaccount.AwsAccount](ctx, lg, nil, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a repository")
	}
	defer srv.Close()

}

func TestAwsAccountServiceGetAll(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awsaccounts-getall.db")

	// insert standard items including teams before this to create joins
	_, err = team.Seed(ctx, lg, cfg, nil)
	if err != nil {
		t.Errorf("unexpected error seeding: [%s]", err.Error())
	}
	insert, err := awsaccount.Seed(ctx, lg, cfg, nil)
	if err != nil {
		t.Errorf("unexpected error seeding: [%s]", err.Error())
	}

	// generate the service useing default
	srv := awsaccount.Default[*awsaccount.AwsAccount](ctx, lg, cfg)
	if srv == nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
	} else {
		defer srv.Close()
	}

	// fetch everything
	res, err := srv.GetAllAccounts()
	if err != nil {
		t.Errorf("unexpected error getting data from service: [%s]", err.Error())
	}
	if len(res) != len(insert) {
		t.Errorf("incorrect number of results found in service")
	}

}
