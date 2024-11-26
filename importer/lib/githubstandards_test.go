package lib

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
)

func Test_processStandards(t *testing.T) {
	var (
		adaptor      dbs.Adaptor
		err          error
		ctx                 = context.Background()
		dir          string = t.TempDir() // "./"
		dbFile       string = filepath.Join(dir, "test.db")
		sourceFile   string = filepath.Join(dir, "standards.json") // "../../convertor/converted/github_standards.json"
		units        []*models.Unit
		teams        []*models.GitHubTeam
		repositories []*models.GitHubRepository
		stds         []*models.GitHubRepositoryStandard
	)
	fakerextras.AddProviders()

	units = fakermany.Fake[*models.Unit](2)
	teams = fakermany.Fake[*models.GitHubTeam](3)
	repositories = fakermany.Fake[*models.GitHubRepository](3)
	stds = fakermany.Fake[*models.GitHubRepositoryStandard](3)

	for _, team := range teams {
		team.Units = units
	}
	for _, repo := range repositories {
		repo.GitHubTeams = teams
	}
	for i, s := range stds {
		s.GitHubRepository = (*models.GitHubRepositoryForeignKey)(repositories[i])
	}
	// dump test data to file
	structs.ToFile(stds, sourceFile)

	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	err = processGithubStandards(ctx, adaptor, sourceFile)
	if err != nil {
		t.Fatalf(err.Error())
	}

}
