package awsuptime

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"
	"opg-reports/report/internal/service/seed"
	"opg-reports/report/internal/utils"
)

func seedDB(ctx context.Context, log *slog.Logger, conf *config.Config) (inserted *seed.SeedAllResults, err error) {
	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	inserted, err = seeder.All(sqc)
	return
}

// Make sure that the handler finds the correct accounts
func TestHandleGetAwsUptimeAll(t *testing.T) {
	var (
		err  error
		dir  = t.TempDir()
		ctx  = t.Context()
		conf = config.NewConfig()
		log  = utils.Logger("DEBUG", "TEXT")
	)
	// overwrite the database location
	conf.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awsuptime-getall.db")
	// capture the inserted data
	inserted, err := seedDB(ctx, log, conf)
	// generate a repository and service
	repository, _ := sqlr.NewWithSelect[*api.AwsUptime](ctx, log, conf)
	service, _ := api.New[*api.AwsUptime](ctx, log, conf)
	// grab the result
	response, err := handleGetAwsUptimeAll(ctx, log, conf, service, repository, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	// make sure all counts match
	if len(inserted.AwsUptime) != response.Body.Count {
		t.Errorf("count doesnt match count of inserted records, expected [%d] actual [%d]", len(inserted.AwsUptime), response.Body.Count)
	}
	if len(inserted.AwsUptime) != len(response.Body.Data) {
		t.Errorf("data length doesnt match count of inserted records, expected [%d] actual [%d]", len(inserted.AwsUptime), len(response.Body.Data))
	}

}
