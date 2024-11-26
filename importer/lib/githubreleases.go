package lib

import (
	"context"
	"fmt"
	"slices"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
)

// processGithubReleases handles importing github standards file with the structure of:
//   - GitHubRelease
//     -- GitHubRepository
//     ---- GitHubTeam
func processGithubReleases(ctx context.Context, adaptor dbs.Adaptor, path string) (res any, err error) {
	var (
		releases   []*models.GitHubRelease
		repos      []*models.GitHubRepository
		teams      []*models.GitHubTeam
		reposFound []string
		teamsFound []string
	)
	if adaptor == nil {
		err = fmt.Errorf("adaptor is nil")
		return
	}
	// read the file and convert into standards list
	if err = structs.UnmarshalFile(path, &releases); err != nil {
		return
	}

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.All()...)
	if err != nil {
		return
	}

	// now get the unique repositories and attach teams
	for _, release := range releases {
		var repo *models.GitHubRepository = (*models.GitHubRepository)(release.GitHubRepository)
		var found *models.GitHubRepository
		var relTeams []*models.GitHubTeam
		// sort out the team
		if !slices.Contains(reposFound, repo.UniqueValue()) {
			reposFound = append(reposFound, repo.UniqueValue())
			repos = append(repos, repo)
		}
		for _, r := range repos {
			if r.UniqueValue() == repo.UniqueValue() {
				found = r
			}
		}
		// now the teams
		for _, team := range repo.GitHubTeams {
			// add to main list of teams
			if !slices.Contains(teamsFound, team.UniqueValue()) {
				teamsFound = append(teamsFound, team.UniqueValue())
				teams = append(teams, team)
			}
			for _, t := range teams {
				if t.UniqueValue() == team.UniqueValue() {
					relTeams = append(relTeams, t)
				}
			}
		}

		if found != nil {
			found.GitHubTeams = relTeams
			release.GitHubRepository = (*models.GitHubRepositoryForeignKey)(found)
		}

	}

	// insert teams
	if _, err = crud.Insert(ctx, adaptor, &models.GitHubTeam{}, teams...); err != nil {
		return
	}
	// insert repos
	if _, err = crud.Insert(ctx, adaptor, &models.GitHubRepository{}, repos...); err != nil {
		return
	}

	// now update the release repo id
	for _, rel := range releases {
		rel.GitHubRepositoryID = rel.GitHubRepository.ID
	}
	// insert releases
	if _, err = crud.Insert(ctx, adaptor, &models.GitHubRelease{}, releases...); err != nil {
		return
	}

	//-- many to many joins

	// repo <-> team
	repoteams := []*models.GitHubRepositoryGitHubTeam{}
	for _, repo := range repos {
		for _, team := range repo.GitHubTeams {
			var join = &models.GitHubRepositoryGitHubTeam{
				GitHubRepositoryID: repo.ID,
				GitHubTeamID:       team.ID,
			}
			repoteams = append(repoteams, join)
		}
	}

	if _, err = crud.Insert(ctx, adaptor, &models.GitHubRepositoryGitHubTeam{}, repoteams...); err != nil {
		return
	}

	// team <-> unit
	teamunits := []*models.GitHubTeamUnit{}
	for _, team := range teams {
		for _, unit := range team.StandardUnits() {
			// insert the unit
			crud.Insert(ctx, adaptor, &models.Unit{}, unit)

			var join = &models.GitHubTeamUnit{
				UnitID:       unit.ID,
				GitHubTeamID: team.ID,
			}
			teamunits = append(teamunits, join)
		}
	}

	if _, err = crud.Insert(ctx, adaptor, &models.GitHubTeamUnit{}, teamunits...); err != nil {
		return
	}

	res = releases
	return

}
