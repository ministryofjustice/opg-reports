package githubcodeowners

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
func TestHandleGetGithubCodeOwnersAll(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = utils.Logger("ERROR", "TEXT")
	)
	// overwrite the database location
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-githubowners-all.db")
	// capture the inserted data
	inserted, err := seedDB(ctx, lg, cfg)
	// generate a repository and service
	repository, _ := sqlr.NewWithSelect[*api.GithubCodeOwner](ctx, lg, cfg)
	service, _ := api.New[*api.GithubCodeOwner](ctx, lg, cfg)
	// grab the result
	response, err := handleGetGithubCodeOwnersAll(ctx, lg, cfg, service, repository, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	// make sure all counts match
	if len(inserted.GithubCodeOwners) != response.Body.Count {
		t.Errorf("count doesnt match count of inserted records, expected [%d] actual [%d]", len(inserted.GithubCodeOwners), response.Body.Count)
	}
	if len(inserted.GithubCodeOwners) != len(response.Body.Data) {
		t.Errorf("data length doesnt match count of inserted records, expected [%d] actual [%d]", len(inserted.GithubCodeOwners), len(response.Body.Data))
	}

}
