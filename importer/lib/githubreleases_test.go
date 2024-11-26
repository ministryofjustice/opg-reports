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

func Test_processGithubReleases(t *testing.T) {
	var (
		adaptor dbs.Adaptor
		err     error
		res     any
		ok      bool
		ctx     context.Context = context.Background()
		// dir        string          = "./"
		// sourceFile string          = "../../collectors/githubreleases/data/2024-11-25_github_releases.json"
		dir        string = t.TempDir()
		sourceFile string = filepath.Join(dir, "data.json")
		dbFile     string = filepath.Join(dir, "test.db")
		repos      []*models.GitHubRepository
		teams      []*models.GitHubTeam
		releases   []*models.GitHubRelease
		result     []*models.GitHubRelease
	)
	// structs.UnmarshalFile(sourceFile, &releases)

	fakerextras.AddProviders()

	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	// make some fake units
	repos = fakermany.Fake[*models.GitHubRepository](3)
	// some fake accounts
	teams = fakermany.Fake[*models.GitHubTeam](3)
	// join them up
	for _, r := range repos {
		r.GitHubTeams = teams
	}
	// some fake uptimes
	releases = fakermany.Fake[*models.GitHubRelease](3)
	// join the accounts and units
	for i, rel := range releases {
		rel.GitHubRepository = (*models.GitHubRepositoryForeignKey)(repos[i])
	}

	structs.ToFile(releases, sourceFile)

	res, err = processGithubReleases(ctx, adaptor, sourceFile)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if result, ok = res.([]*models.GitHubRelease); !ok {
		t.Errorf("failed to change result to type")
	}

	if len(releases) != len(result) {
		t.Errorf("number of returned results dont match originals - expected [%d] actual [%v]", len(releases), len(result))
	}

}
