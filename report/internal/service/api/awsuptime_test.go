package api

import (
	"path/filepath"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
)

// TestApiServiceGetAllAwsUptime uses get all func - which is generally bad idea with real data
func TestApiServiceGetAllAwsUptime(t *testing.T) {

	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./api-get-all-awsuptime.db")
	inserted, err := seedDB(ctx, log, conf)

	store := sqlr.DefaultWithSelect[*AwsUptime](ctx, log, conf)
	service := Default[*AwsUptime](ctx, log, conf)

	data, err := service.GetAllAwsUptime(store)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(inserted.AwsUptime) != len(data) {
		t.Errorf("mismatching number of records: expected [%d] actual [%v]", len(inserted.AwsUptime), len(data))
	}

}
