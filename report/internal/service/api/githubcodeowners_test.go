package api

import (
	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"
	"path/filepath"
	"testing"
)

// TestApiServiceGetAllGithubCodeOwners uses get all func - which is generally
// bad idea with real data
func TestApiServiceGetAllGithubCodeOwners(t *testing.T) {

	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("ERROR", "TEXT")
	)
	// set config values
	conf.Database.Path = filepath.Join(dir, "./api-get-all-codeowners.db")
	inserted, err := seedDB(ctx, log, conf)

	store := sqlr.DefaultWithSelect[*GithubCodeOwner](ctx, log, conf)
	service := Default[*GithubCodeOwner](ctx, log, conf)

	data, err := service.GetAllGithubCodeOwners(store)

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(inserted.GithubCodeOwners) != len(data) {
		t.Errorf("mismatching number of records: expected [%d] actual [%v]", len(inserted.AwsUptime), len(data))
	}

}
