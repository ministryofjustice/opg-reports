package awscost

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestGetGroupedCostsOptions(t *testing.T) {
	var (
		err   error
		valid bool
		stmt  *sqlr.BoundStatement
		data  *sqlParams
		dir   = t.TempDir()
		ctx   = t.Context()
		conf  = config.NewConfig()
		log   = utils.Logger("ERROR", "TEXT")
	)
	conf.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awscost-getopts.db")

	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	seeder.Teams(sqc)
	seeder.AwsAccounts(sqc)
	seeder.AwsCosts(sqc)

	repo, err := sqlr.NewWithSelect[*AwsCost](ctx, log, conf)

	opts := &GetGroupedCostsOptions{
		StartDate:   "2024-01",
		EndDate:     "2024-03",
		DateFormat:  "%Y-%m",
		Team:        "",
		Service:     "",
		Region:      "",
		Account:     "",
		Environment: "",
	}
	// Test empty general grouping
	stmt, data = opts.Statement()
	valid, _, err = repo.ValidateSelect(stmt)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !valid {
		t.Errorf("unexpected validation failure")
	}
	if data == nil {
		t.Errorf("incorrect data items")
	}

	// Test team == true, so should be in the select, but not in the data
	opts.Team = "true"
	stmt, data = opts.Statement()
	valid, _, err = repo.ValidateSelect(stmt)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !valid {
		t.Errorf("unexpected validation failure")
	}
	if data.Team != "" {
		t.Errorf("team should not be set as a filter")
	}

	// Test team filtering by a value is working
	opts.Team = "Team0A"
	stmt, data = opts.Statement()
	valid, _, err = repo.ValidateSelect(stmt)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !valid {
		t.Errorf("unexpected validation failure")
	}
	if data.Team == "" {
		t.Errorf("team should be set as a filter")
	}

	// Test Team filter & account grouping
	opts.Team = "Team0A"
	opts.Account = "true"
	stmt, data = opts.Statement()
	valid, _, err = repo.ValidateSelect(stmt)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !valid {
		t.Errorf("unexpected validation failure")
	}
	if data.Team == "" {
		t.Errorf("team should be set as a filter")
	}
	if data.Account != "" {
		t.Errorf("aws_account_id should NOT set as a filter")
	}

	// Test Team and account filter
	opts.Team = "Team0A"
	opts.Account = "Acc01"
	stmt, data = opts.Statement()
	valid, _, err = repo.ValidateSelect(stmt)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !valid {
		t.Errorf("unexpected validation failure")
	}
	if data.Team == "" {
		t.Errorf("team should be set as a filter")
	}
	if data.Account == "" {
		t.Errorf("aws_account_id should be set as a filter")
	}

	// Test grouping by most - filtering by env
	opts.Team = "true"
	opts.Account = "true"
	opts.Environment = "production"
	opts.Region = "true"
	opts.Service = "true"
	stmt, data = opts.Statement()
	valid, _, err = repo.ValidateSelect(stmt)
	if data.Environment == "" {
		t.Errorf("environment should be set as a filter")
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if !valid {
		t.Errorf("unexpected validation failure")
	}
	if data.Team != "" {
		t.Errorf("team should not be set as a filter")
	}
	if data.Account != "" {
		t.Errorf("aws_account_id should not be set as a filter")
	}

	if data.Region != "" {
		t.Errorf("region should not be set as a filter")
	}
	if data.Service != "" {
		t.Errorf("service should not be set as a filter")
	}

}
