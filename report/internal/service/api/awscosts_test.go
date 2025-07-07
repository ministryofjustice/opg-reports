package api

import (
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
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
	_, _, inserted := seedDB(ctx, log, conf)

	store := sqlr.DefaultWithSelect[*AwsCost](ctx, log, conf)
	service := Default[*AwsCost](ctx, log, conf)

	data, err := service.GetAllAwsCosts(store)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(inserted) != len(data) {
		t.Errorf("mismatching number of records: expected [%d] actual [%v]", len(inserted), len(data))
	}

}
