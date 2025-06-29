package awscost

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestGetGroupedCostsOptions(t *testing.T) {
	var (
		err   error
		valid bool
		stmt  *sqldb.BoundStatement
		data  *sqlParams
		dir   = t.TempDir()
		ctx   = t.Context()
		cfg   = config.NewConfig()
		lg    = utils.Logger("ERROR", "TEXT")
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awscost-getopts.db")
	// seed teams, accounts & costs to get some dummy data
	team.Seed(ctx, lg, cfg, nil)
	awsaccount.Seed(ctx, lg, cfg, nil)
	Seed(ctx, lg, cfg, nil)
	repo, err := sqldb.New[*AwsCost](ctx, lg, cfg)

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
