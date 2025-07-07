package awsaccounts

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/api"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func seedDB(ctx context.Context, log *slog.Logger, conf *config.Config) (inserted []*sqlr.BoundStatement) {
	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	seeder.Teams(sqc)
	inserted, _ = seeder.AwsAccounts(sqc)
	return
}

// Make sure that the handler finds the correct accounts
func TestHandleGetAwsAccountsAll(t *testing.T) {
	var (
		err  error
		dir  = t.TempDir()
		ctx  = t.Context()
		conf = config.NewConfig()
		log  = utils.Logger("DEBUG", "TEXT")
	)
	// overwrite the database location
	conf.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awsaccounts-getall.db")
	// capture the inserted data
	inserted := seedDB(ctx, log, conf)
	// generate a repository and service
	repository, _ := sqlr.NewWithSelect[*api.AwsAccount](ctx, log, conf)
	service, _ := api.New[*api.AwsAccount](ctx, log, conf)
	// grab the result
	response, err := handleGetAwsAccountsAll(ctx, log, conf, service, repository, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	utils.Debug(response)
	// make sure all counts match
	if len(inserted) != response.Body.Count {
		t.Errorf("count doesnt match count of inserted records, expected [%d] actual [%d]", len(inserted), response.Body.Count)
	}
	if len(inserted) != len(response.Body.Data) {
		t.Errorf("data length doesnt match count of inserted records, expected [%d] actual [%d]", len(inserted), len(response.Body.Data))
	}

	// now test that the data returned has the correct account ids
	found := 0
	for _, item := range response.Body.Data {
		for _, insert := range inserted {
			if insert.Returned.(string) == item.ID {
				found++
			}
		}
	}
	if found != len(inserted) {
		t.Errorf("dit not find all inserted records from in the response")
	}
}
