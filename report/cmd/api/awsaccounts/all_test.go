package awsaccounts

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/service/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func seed(ctx context.Context, lg *slog.Logger, cfg *config.Config) (inserted []*sqlr.BoundStatement) {
	team.Seed(ctx, lg, cfg, nil)
	inserted, _ = awsaccount.Seed(ctx, lg, cfg, nil)
	return
}

// Make sure that the handler finds the correct accounts
func TestHandleGetAwsAccountsAll(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)
	// overwrite the database location
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awsaccounts-getall.db")
	// capture the inserted data
	inserted := seed(ctx, lg, cfg)
	// generate a repository and service
	repository, _ := sqlr.New[*AwsAccount](ctx, lg, cfg)
	service, _ := awsaccount.NewService(ctx, lg, cfg, repository)
	// grab the result
	response, err := handleGetAwsAccountsAll(ctx, lg, cfg, service, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
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
