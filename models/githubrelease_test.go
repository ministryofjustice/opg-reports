package models_test

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
)

// Interface checks
var (
	_ dbs.Table           = &models.GitHubRelease{}
	_ dbs.CreateableTable = &models.GitHubRelease{}
	_ dbs.Insertable      = &models.GitHubRelease{}
	_ dbs.Row             = &models.GitHubRelease{}
	_ dbs.InsertableRow   = &models.GitHubRelease{}
	_ dbs.Record          = &models.GitHubRelease{}
)

// TestModelsGitHubReleaseCRUD checks the github team table creation
// and inserting series of fake records works as expected
func TestModelsGitHubReleaseCRUD(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err     error
		adaptor *adaptors.Sqlite
		n       int                     = 100
		ctx     context.Context         = context.Background()
		dir     string                  = t.TempDir()
		path    string                  = filepath.Join(dir, "test.db")
		units   []*models.GitHubRelease = fakermany.Fake[*models.GitHubRelease](n)
		results []*models.GitHubRelease
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &models.GitHubRelease{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}
	_, err = crud.CreateIndexes(ctx, adaptor, &models.GitHubRelease{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &models.GitHubRelease{}, units...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(results) != len(units) {
		t.Errorf("created records do not match expacted number - [%d] actual [%v]", len(units), len(results))
	}

}

var selectReleases = `
SELECT
	github_releases.*,
	json_object(
		'id', github_repositories.id,
		'full_name', github_repositories.full_name,
		'owner', github_repositories.owner
	) as github_repository
FROM github_releases
LEFT JOIN github_repositories on github_repositories.id = github_releases.github_repository_id
GROUP BY github_releases.id
ORDER BY github_releases.date ASC;
`

func TestModelsGithubReleaseRepoJoin(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err      error
		adaptor  *adaptors.Sqlite
		ctx      context.Context = context.Background()
		dir      string          = t.TempDir()
		path     string          = filepath.Join(dir, "test.db")
		releases []*models.GitHubRelease
		repos    []*models.GitHubRepository
		results  []*models.GitHubRelease
	)
	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	defer adaptor.DB().Close()

	// create test repos
	repos, err = testDBbuilder(ctx, adaptor, &models.GitHubRepository{}, fakermany.Fake[*models.GitHubRepository](5))
	if err != nil {
		t.Fatalf(err.Error())
	}

	// create test releaases
	releases = fakermany.Fake[*models.GitHubRelease](10)
	// now add a repo to each release
	for _, rel := range releases {
		var repo = fakerextras.Choice(repos)
		var join = models.GitHubRepositoryForeignKey(*repo)
		rel.GitHubRepositoryID = repo.ID
		rel.GitHubRepository = &join
	}
	// now save the items
	_, err = testDBbuilder(ctx, adaptor, &models.GitHubRelease{}, releases)
	if err != nil {
		t.Fatalf(err.Error())
	}

	results, err = crud.Select[*models.GitHubRelease](ctx, adaptor, selectReleases, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	//  check length matches
	if len(releases) != len(results) {
		t.Errorf("length mismatch - expected [%d] actual [%v]", len(releases), len(results))
	}

	// now check the selected results match the generated version
	for _, res := range results {
		var release *models.GitHubRelease
		for _, r := range releases {
			if r.ID == res.ID {
				release = r
			}
		}
		if release == nil {
			t.Errorf("failed to find release")
		}
		if release.GitHubRepositoryID != res.GitHubRepositoryID {
			t.Errorf("repo ID mismatch")
		}
		if release.GitHubRepository.ID != res.GitHubRepository.ID {
			t.Errorf("repo ID mismatch")
		}
	}

}
