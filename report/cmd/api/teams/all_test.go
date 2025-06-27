package teams

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func seed(ctx context.Context, lg *slog.Logger, cfg *config.Config) (inserted []*sqldb.BoundStatement) {
	inserted, _ = team.Seed(ctx, lg, cfg, nil)
	return
}

// Make sure that the handler finds the correct accounts
func TestHandleGetTeamsAll(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)
	// overwrite the database location
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-teams-all.db")
	// capture the inserted data
	inserted := seed(ctx, lg, cfg)
	// generate a repository and service
	repository, _ := sqldb.New[*Team](ctx, lg, cfg)
	service, _ := team.NewService(ctx, lg, cfg, repository)
	// grab the result
	response, err := handleGetTeamsAll(ctx, lg, cfg, service, nil)
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
	//  - casting to int64 to compare
	found := 0
	for _, item := range response.Body.Data {
		for _, insert := range inserted {
			if insert.Returned.(int64) == (int64)(item.ID) {
				found++
			}
		}
	}
	if found != len(inserted) {
		t.Errorf("dit not find all inserted records from in the response [%d] [%d]", len(inserted), found)
	}
}
