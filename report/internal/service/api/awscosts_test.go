package api

import (
	"path/filepath"
	"testing"
	"time"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

// TestApiServiceGetAllAwsCosts uses get all func - which is generally bad idea with real
// data (which has over 2m records)
func TestApiServiceGetAllAwsCosts(t *testing.T) {

	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./api-get-all-awscosts.db")
	inserted, err := seedDB(ctx, log, conf)

	store := sqlr.DefaultWithSelect[*AwsCost](ctx, log, conf)
	service := Default[*AwsCost](ctx, log, conf)

	data, err := service.GetAllAwsCosts(store)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(inserted.AwsCosts) != len(data) {
		t.Errorf("mismatching number of records: expected [%d] actual [%v]", len(inserted.AwsCosts), len(data))
	}

}

func TestApiServicePutAwsCosts(t *testing.T) {

	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
		data        = []*AwsCost{
			{Cost: "1.152", Date: "2025-05-31", Region: "eu-west-1", Service: "Amazon Virtual Private Cloud", AwsAccountID: "004B"},
		}
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./api-put-awscosts.db")
	seedDB(ctx, log, conf)

	store := sqlr.DefaultWithSelect[*AwsCost](ctx, log, conf)
	service := Default[*AwsCost](ctx, log, conf)

	inserted, err := service.PutAwsCosts(store, data)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(inserted) != len(data) {
		t.Errorf("mismatching number of records: expected [%d] actual [%v]", len(inserted), len(data))
	}

}

func TestApiServiceGetGroupedAwsCosts(t *testing.T) {
	var (
		err   error
		end   time.Time = time.Now().UTC()
		start time.Time = end.AddDate(0, -2, 0)
		dir   string    = t.TempDir()
		ctx             = t.Context()
		conf            = config.NewConfig()
		log             = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./api-get-grouped-awscosts.db")
	// seed database
	seedDB(ctx, log, conf)

	store := sqlr.DefaultWithSelect[*AwsCostGrouped](ctx, log, conf)
	service := Default[*AwsCostGrouped](ctx, log, conf)

	opts := &GetAwsCostsGroupedOptions{
		StartDate:  start.Format(utils.DATE_FORMATS.YMD),
		EndDate:    end.Format(utils.DATE_FORMATS.YMD),
		DateFormat: utils.GRANULARITY_TO_FORMAT["month"],
	}
	// this should return jsut 1 result, as should all be merged by month
	data, err := service.GetGroupedAwsCosts(store, opts)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(data) != 1 {
		t.Errorf("expected 1, actual: %d", len(data))
	}
	// grouping by account & team
	opts.Team = "true"
	data, err = service.GetGroupedAwsCosts(store, opts)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(data) != 4 {
		t.Errorf("expected 4, actual: %d", len(data))
	}

	// test grouping by account - so should be month*num_of_accounts
	opts.Team = ""
	opts.Account = "true"
	data, err = service.GetGroupedAwsCosts(store, opts)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(data) != 5 {
		t.Errorf("expected 5 result, actual: %d", len(data))
	}
	// grouping by account & team should be 5
	opts.Account = "true"
	opts.Team = "true"
	data, err = service.GetGroupedAwsCosts(store, opts)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(data) != 5 {
		t.Errorf("expected 5 result, actual: %d", len(data))
	}
	// check filtering by a team - should only return 1 team data
	opts.Account = ""
	opts.Team = "TEAM-A"
	data, err = service.GetGroupedAwsCosts(store, opts)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(data) != 1 {
		t.Errorf("expected 1, actual: %d", len(data))
	}
	if data[0].TeamName != string(opts.Team) {
		t.Errorf("team found does not match filter")
	}

}
