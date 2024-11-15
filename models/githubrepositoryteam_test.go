package models_test

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/internal/pretty"
	"github.com/ministryofjustice/opg-reports/models"
)

// Interface checks
var (
	_ dbs.Table           = &models.GitHubRepositoryGitHubTeam{}
	_ dbs.CreateableTable = &models.GitHubRepositoryGitHubTeam{}
	_ dbs.Insertable      = &models.GitHubRepositoryGitHubTeam{}
	_ dbs.Row             = &models.GitHubRepositoryGitHubTeam{}
	_ dbs.InsertableRow   = &models.GitHubRepositoryGitHubTeam{}
	_ dbs.Record          = &models.GitHubRepositoryGitHubTeam{}
)

var selectRepos string = `
SELECT
	github_repositories.*,
	json_group_array(json_object('id', github_teams.id,'name', github_teams.name)) as github_teams
FROM github_repositories
LEFT JOIN github_repositories_github_teams on github_repositories_github_teams.github_repository_id = github_repositories.id
LEFT JOIN github_teams on github_teams.id = github_repositories_github_teams.github_team_id
GROUP BY github_repositories.id
ORDER BY github_repositories.full_name ASC;
`

// TestModelsGithubTeamUnitJoin checks the join logic from
// repo->teams is working
func TestModelsGithubRepoTeamJoin(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
		repos   []*models.GitHubRepository
		results []*models.GitHubRepository
		teams   []*models.GitHubTeam
		joins   []*models.GitHubRepositoryGitHubTeam
	)
	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	defer adaptor.DB().Close()

	repos, err = testDBbuilder(ctx, adaptor, &models.GitHubRepository{}, fakermany.Fake[*models.GitHubRepository](10))
	if err != nil {
		t.Fatalf(err.Error())
	}
	teams, err = testDBbuilder(ctx, adaptor, &models.GitHubTeam{}, fakermany.Fake[*models.GitHubTeam](4))
	if err != nil {
		t.Fatalf(err.Error())
	}

	// now we create the joins and insert them
	for _, repo := range repos {
		var list = []*models.GitHubTeam{}
		var picked = fakerextras.Choices(teams, 2)
		// dont add duplicates
		for _, i := range picked {
			if !slices.Contains(list, i) {
				list = append(list, i)
			}
		}
		// set the team on the unit
		repo.GitHubTeams = list
		for _, gt := range list {
			joins = append(joins, &models.GitHubRepositoryGitHubTeam{GitHubRepositoryID: repo.ID, GitHubTeamID: gt.ID})
		}
	}
	// insert the joins and create the table etc
	_, err = testDBbuilder(ctx, adaptor, &models.GitHubRepositoryGitHubTeam{}, joins)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// now select the units to and see if the teams are included!
	results, err = crud.Select[*models.GitHubRepository](ctx, adaptor, selectRepos, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	pretty.Print(results)
	// now check the results contain the correct join data
	for _, result := range results {
		// grab actually returned teams
		var actualTeams = result.GitHubTeams
		// get the teams that were generated for this unit
		var expectedTeams = []*models.GitHubTeam{}
		for _, repo := range repos {
			if repo.ID == result.ID {
				expectedTeams = repo.GitHubTeams
			}
		}
		// if the counts dont match, throw an error
		if len(expectedTeams) != len(actualTeams) {
			t.Errorf("actual teams do not match expected versions:")
			fmt.Println("expected:")
			pretty.Print(expectedTeams)
			fmt.Println("actual:")
			pretty.Print(actualTeams)
		}
		// now compare both sides to make sure content match
		for _, exp := range expectedTeams {
			var found = false
			for _, act := range actualTeams {
				if act.Name == exp.Name {
					found = true
				}
			}
			if !found {
				t.Errorf("set teams dont match returned teams")
			}
		}

	}

}
