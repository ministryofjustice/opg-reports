package lib

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/seed"
)

// SeedDB is called when there is currently no database at `dbPath` with
// the intention to create a new version, but with dummy data.
//
// Creates a series of related models and then calls methods in `seed`
// to generate the database and create joins etc correctly.
//
// If there is an error in this function a panic is called as a lack of
// database means the api cannot function.
func SeedDB(ctx context.Context, dbPath string) {
	var (
		err     error
		adaptor dbs.Adaptor
	)
	// use sqlite
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[seed] error with adaptor", slog.String("err", err.Error()))
		panic(err)
	}
	// bootstrap will generate db tables and indexes from the models
	err = crud.Bootstrap(ctx, adaptor, models.All()...)
	if err != nil {
		slog.Error("[seed] error with bootstrap", slog.String("err", err.Error()))
		panic(err)
	}

	units := getUnits()
	// github
	teams := getGitHubTeams(units)
	repositories := getGitHubRepositories(teams)
	standards := getGitHubStandards(repositories)
	releases := getGitHubReleaes(repositories)
	// aws
	accounts := getAwsAccounts(units)
	costs := getAwsCosts(accounts)
	uptime := getAwsUptime(accounts)

	// seed standards
	_, err = seed.GitHubStandards(ctx, adaptor, standards)
	if err != nil {
		slog.Error("[seed] error with seeding github standards", slog.String("err", err.Error()))
		panic(err)
	}
	// seed releases
	_, err = seed.GitHubReleases(ctx, adaptor, releases)
	if err != nil {
		slog.Error("[seed] error with seed github releases", slog.String("err", err.Error()))
		panic(err)
	}
	// seed costs
	_, err = seed.AwsCosts(ctx, adaptor, costs)
	if err != nil {
		slog.Error("[seed] error with seeding aws costs", slog.String("err", err.Error()))
		panic(err)
	}
	// seed uptime
	_, err = seed.AwsUptime(ctx, adaptor, uptime)
	if err != nil {
		slog.Error("[seed] error with seeding aws uptime", slog.String("err", err.Error()))
		panic(err)
	}

}

func getAwsUptime(accounts []*models.AwsAccount) (uptime []*models.AwsUptime) {
	uptime = fakermany.Fake[*models.AwsUptime](14400)
	for _, up := range uptime {
		var account = fakerextras.Choice(accounts)
		up.AwsAccount = (*models.AwsAccountForeignKey)(account)
		up.Unit = account.Unit
	}
	return
}

func getAwsCosts(accounts []*models.AwsAccount) (costs []*models.AwsCost) {
	costs = fakermany.Fake[*models.AwsCost](100000)
	for _, cost := range costs {
		var account = fakerextras.Choice(accounts)
		cost.AwsAccount = (*models.AwsAccountForeignKey)(account)
		cost.Unit = account.Unit
	}
	return
}

func getAwsAccounts(units []*models.Unit) (accounts []*models.AwsAccount) {
	accounts = fakermany.Fake[*models.AwsAccount](5)
	for _, account := range accounts {
		var unit = fakerextras.Choice(units)
		account.Unit = (*models.UnitForeignKey)(unit)
	}
	return
}

func getGitHubReleaes(repositories []*models.GitHubRepository) (releases []*models.GitHubRelease) {
	releases = fakermany.Fake[*models.GitHubRelease](1000)
	for _, release := range releases {
		var repo = fakerextras.Choice(repositories)
		release.GitHubRepository = (*models.GitHubRepositoryForeignKey)(repo)

	}
	return
}

func getGitHubStandards(repositories []*models.GitHubRepository) (standards []*models.GitHubRepositoryStandard) {
	var i = len(repositories)
	standards = fakermany.Fake[*models.GitHubRepositoryStandard](i)
	for i, s := range standards {
		s.GitHubRepository = (*models.GitHubRepositoryForeignKey)(repositories[i])
	}
	return
}

func getGitHubRepositories(teams []*models.GitHubTeam) (repositories []*models.GitHubRepository) {
	repositories = fakermany.Fake[*models.GitHubRepository](50)
	for _, repo := range repositories {
		repo.GitHubTeams = fakerextras.Choices(teams, 1)
	}
	return
}

func getGitHubTeams(units []*models.Unit) (teams []*models.GitHubTeam) {
	teams = fakermany.Fake[*models.GitHubTeam](6)
	for _, team := range teams {
		team.Units = fakerextras.Choices(units, 1)
	}
	return
}

func getUnits() (units []*models.Unit) {
	units = fakermany.Fake[*models.Unit](5)
	return
}
