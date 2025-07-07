package api

import (
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestApiServiceGetAllAwsAccounts(t *testing.T) {

	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./api-get-all-awsaccounts.db")
	_, inserted, _ := seedDB(ctx, log, conf)

	store := sqlr.DefaultWithSelect[*AwsAccount](ctx, log, conf)
	service := Default[*AwsAccount](ctx, log, conf)

	data, err := service.GetAllAwsAccounts(store)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(inserted) != len(data) {
		t.Errorf("mismatching number of records: expected [%d] actual [%v]", len(inserted), len(data))
	}

}
