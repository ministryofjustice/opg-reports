package awscosts

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
	seeder.AwsAccounts(sqc)
	inserted, _ = seeder.AwsCosts(sqc)

	return
}

// Make sure that the handler finds the correct accounts
func TestHandleGetAwsCostsTop20(t *testing.T) {
	var (
		err  error
		dir  = t.TempDir()
		ctx  = t.Context()
		conf = config.NewConfig()
		log  = utils.Logger("ERROR", "TEXT")
	)
	// overwrite the database location
	conf.Database.Path = fmt.Sprintf("%s/%s", dir, "test-costs-top20.db")
	// capture the inserted data
	inserted := seedDB(ctx, log, conf)
	// generate a repository and service
	repository, _ := sqlr.NewWithSelect[*api.AwsCost](ctx, log, conf)
	service, _ := api.New[*api.AwsCost](ctx, log, conf)
	// grab the result
	response, err := handleGetAwsCostsTop20(ctx, log, conf, service, repository, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if response.Body.Count != 20 {
		t.Errorf("expected 20 results to be returned :%v", response.Body.Count)
	}
	// make sure all counts match
	if len(inserted) <= response.Body.Count {
		t.Errorf("count mismatch")
	}

}
