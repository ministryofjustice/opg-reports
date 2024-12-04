package handlers_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/seed"
	"github.com/ministryofjustice/opg-reports/servers/api/handlers"
	"github.com/ministryofjustice/opg-reports/servers/api/lib"
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

// TestApiHandlersGitHubReleasesListHandler creates and then seeds a dummy database
// containing releases and join data (teams / units).
// After creating the seeded db it calls the api end point handler directly using
// configured inputs and checks that the data returned aligns with the created data.
//
// Checks the correctness of the sql statement used in the api handler
func TestApiHandlersGitHubReleasesListHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *inout.GitHubReleasesListResponse
		dir      string = t.TempDir()
		// dir      string                     = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		repos    []*models.GitHubRepository
		units    []*models.Unit
		teams    []*models.GitHubTeam
		releases []*models.GitHubRelease
		// inserted []*models.GitHubTeam = []*models.GitHubTeam{}
	)
	fakerextras.AddProviders()

	units = fakermany.Fake[*models.Unit](5)
	teams = fakermany.Fake[*models.GitHubTeam](5)
	repos = fakermany.Fake[*models.GitHubRepository](5)
	releases = fakermany.Fake[*models.GitHubRelease](5)

	for i, team := range teams {
		var r = repos[i]
		var set = []*models.GitHubRepository{r}
		team.Units = fakerextras.Choose(units, 2)
		team.GitHubRepositories = set
		r.GitHubTeams = []*models.GitHubTeam{team}
	}
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// seed the teams and units
	_, err = seed.GitHubTeams(ctx, adaptor, teams)
	if err != nil {
		t.Fatalf(err.Error())
	}

	for i, release := range releases {
		release.GitHubRepository = (*models.GitHubRepositoryForeignKey)(repos[i])
	}

	// for _, rel := range releases {
	// 	fmt.Printf(">%+v\n", rel)
	// 	fmt.Printf("  >%+v\n", rel.GitHubRepository)
	// 	fmt.Printf("    >%+v\n", rel.GitHubRepository.GitHubTeams)
	// }

	_, err = seed.GitHubReleases(ctx, adaptor, releases)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// should return everything
	in := &inout.DateRangeUnitInput{
		Version:   "v1",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
	}
	response, err = handlers.ApiGitHubReleasesListHandler(ctx, in)

	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.GitHubReleasesListOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.GitHubReleasesListOperationID, response.Body.Operation)
	}
	// check the number of results
	if len(releases) != len(response.Body.Result) {
		t.Errorf("error with number of results - expected [%d] actual [%v]", len(releases), len(response.Body.Result))
	}

}

// TestApiHandlersGitHubReleasesCountHandler generates a seeded database containing
// releases and joined data (teams, units, repos).
// The api handler is then called directly to test that the results align with the
// seeded data correctly.
//
// Checks the sql statement and input parameters.
func TestApiHandlersGitHubReleasesCountHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *inout.GitHubReleasesCountResponse
		dir      string = t.TempDir()
		// dir      string          = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		repos    []*models.GitHubRepository
		units    []*models.Unit
		teams    []*models.GitHubTeam
		releases []*models.GitHubRelease
	)
	fakerextras.AddProviders()

	units = fakermany.Fake[*models.Unit](5)
	teams = fakermany.Fake[*models.GitHubTeam](5)
	repos = fakermany.Fake[*models.GitHubRepository](5)
	releases = fakermany.Fake[*models.GitHubRelease](5)

	for i, team := range teams {
		var r = repos[i]
		var set = []*models.GitHubRepository{r}
		team.Units = fakerextras.Choose(units, 2)
		team.GitHubRepositories = set
		r.GitHubTeams = []*models.GitHubTeam{team}
	}
	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// seed the teams and units
	_, err = seed.GitHubTeams(ctx, adaptor, teams)
	if err != nil {
		t.Fatalf(err.Error())
	}

	for i, release := range releases {
		release.GitHubRepository = (*models.GitHubRepositoryForeignKey)(repos[i])
	}
	// seed releases
	_, err = seed.GitHubReleases(ctx, adaptor, releases)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// should return everything
	in := &inout.RequiredGroupedDateRangeUnitInput{
		Version:   "v1",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
		Interval:  "month",
	}
	in.Resolve(nil)
	response, err = handlers.ApiGitHubReleasesCountHandler(ctx, in)

	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.GitHubReleasesCountOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.GitHubReleasesCountOperationID, response.Body.Operation)
	}

	total := 0
	for _, row := range response.Body.Result {
		total += row.Count
	}
	if len(releases) != total {
		t.Errorf("error with number of results - expected at least [%d] actual [%v]", len(releases), total)
	}

}

// TestApiHandlersGitHubReleasesCountPerUnitHandler creates and seed a database
// containing 100 releases, all contected to the same repository and single team.
// This allows check the grouping and count per time period is working correctly
// as the sum of all counts == number of inserted releases.
func TestApiHandlersGitHubReleasesCountPerUnitHandler(t *testing.T) {
	var (
		err      error
		adaptor  dbs.Adaptor
		response *inout.GitHubReleasesCountPerUnitResponse
		dir      string = t.TempDir()
		// dir      string          = "./"
		dbFile   string          = filepath.Join(dir, "test.db")
		ctxKey   string          = lib.CTX_DB_KEY
		ctx      context.Context = context.WithValue(context.Background(), ctxKey, dbFile)
		repos    []*models.GitHubRepository
		teams    []*models.GitHubTeam
		releases []*models.GitHubRelease
	)
	fakerextras.AddProviders()

	repos = fakermany.Fake[*models.GitHubRepository](1)
	teams = fakermany.Fake[*models.GitHubTeam](1)
	releases = fakermany.Fake[*models.GitHubRelease](50)

	for _, repo := range repos {
		repo.GitHubTeams = teams
	}

	for _, release := range releases {
		release.GitHubRepository = (*models.GitHubRepositoryForeignKey)(repos[0])
	}

	// generate adaptor
	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// seed releases
	_, err = seed.GitHubReleases(ctx, adaptor, releases)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// -- run test

	// should return everything
	in := &inout.RequiredGroupedDateRangeInput{
		Version:   "v1",
		StartDate: fakerextras.TimeStringMin.AddDate(0, 0, -1).Format(dateformats.YMD),
		EndDate:   fakerextras.TimeStringMax.AddDate(0, 0, 1).Format(dateformats.YMD),
		Interval:  "month",
	}
	in.Resolve(nil)
	response, err = handlers.ApiGitHubReleasesCountPerUnitHandler(ctx, in)

	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	// check the response info
	if handlers.GitHubReleasesCountPerUnitOperationID != response.Body.Operation {
		t.Errorf("operation did not match - expected [%s] actual [%v]", handlers.GitHubReleasesCountPerUnitOperationID, response.Body.Operation)
	}
	// as we've inserted only one unit & one repo then the releases should tally exactly
	total := 0
	for _, row := range response.Body.Result {
		total += row.Count
	}
	if len(releases) != total {
		t.Errorf("error with number of results - expected [%d] actual [%v]", len(releases), total)
	}

}
