package awscosts

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/awscost"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func seed(ctx context.Context, lg *slog.Logger, cfg *config.Config) (inserted []*sqldb.BoundStatement) {
	team.Seed(ctx, lg, cfg, nil)
	awsaccount.Seed(ctx, lg, cfg, nil)
	inserted, _ = awscost.Seed(ctx, lg, cfg, nil)
	return
}

// Make sure that the handler finds the correct accounts
func TestHandleGetAwsCostsTop20(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("WARN", "TEXT")
	)
	// overwrite the database location
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-costs-top20.db")
	// capture the inserted data
	inserted := seed(ctx, lg, cfg)
	// generate a repository and service
	repository, _ := sqldb.New[*AwsCost](ctx, lg, cfg)
	service, _ := awscost.NewService(ctx, lg, cfg, repository)
	// grab the result
	response, err := handleGetAwsCostsTop20(ctx, lg, cfg, service, nil)
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
