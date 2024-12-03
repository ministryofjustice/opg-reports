package handlers_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/seed"
	"github.com/ministryofjustice/opg-reports/servers/api/handlers"
	"github.com/ministryofjustice/opg-reports/servers/api/lib"
	"github.com/ministryofjustice/opg-reports/servers/inputs"
)

func TestApiHandlersGitHubTeamsHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *handlers.GitHubTeamsResponse
		dir      string = t.TempDir()
		// dir      string                     = "./"
		dbFile   string                     = filepath.Join(dir, "test.db")
		ctxKey   string                     = lib.CTX_DB_KEY
		ctx      context.Context            = context.WithValue(context.Background(), ctxKey, dbFile)
		units    []*models.Unit             = []*models.Unit{}
		repos    []*models.GitHubRepository = []*models.GitHubRepository{}
		teams    []*models.GitHubTeam       = []*models.GitHubTeam{}
		inserted []*models.GitHubTeam       = []*models.GitHubTeam{}
	)

	units = fakermany.Fake[*models.Unit](5)
	repos = fakermany.Fake[*models.GitHubRepository](10)
	teams = fakermany.Fake[*models.GitHubTeam](6)

	for _, team := range teams {
		team.Units = fakerextras.Choose(units, 2)
		team.GitHubRepositories = fakerextras.Choose(repos, 2)
	}
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}
	inserted, err = seed.GitHubTeams(ctx, adaptor, teams)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// check lengths
	if len(teams) != len(inserted) {
		t.Errorf("error inserting - expected [%d] actual [%v]", len(teams), len(inserted))
	}
	// should return everything
	response, err = handlers.ApiGitHubTeamsListHandler(ctx, &inputs.VersionUnitInput{
		Version: "v1",
		// Unit:    teams[0].Units[0].Name,
	})

	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	// check the response info
	if handlers.GitHubTeamsOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.GitHubTeamsOperationID, response.Body.Operation)
	}
	// check the number of results
	if len(teams) != len(response.Body.Result) {
		t.Errorf("error with number of results - expected [%d] actual [%v]", len(teams), len(response.Body.Result))
	}

	// check for a particular unit - grab the first one in the results
	unit := response.Body.Result[0].Units[0].Name
	expected := 0
	for _, r := range response.Body.Result {
		for _, u := range r.Units {
			if u.Name == unit {
				expected += 1
			}
		}
	}
	response, err = handlers.ApiGitHubTeamsListHandler(ctx, &inputs.VersionUnitInput{
		Version: "v1",
		Unit:    unit,
	})
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	// check the number of results
	if expected != len(response.Body.Result) {
		t.Errorf("error with number of results for unit filter - expected [%d] actual [%v]", expected, len(response.Body.Result))
	}

}
