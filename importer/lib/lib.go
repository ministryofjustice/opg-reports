package lib

import (
	"context"
	"slices"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
)

// Arguments represents all the named arguments for this collector
type Arguments struct {
	DatabasePath string
	SourceFile   string
	Type         string
}

type importerFunc func(ctx context.Context, adaptor dbs.Adaptor, path string) error

var typeProcessors = map[string]importerFunc{
	"standards": processStandards,
}

func processStandards(ctx context.Context, adaptor dbs.Adaptor, path string) (err error) {
	var (
		standards    []*models.GitHubRepositoryStandard
		repositories []*models.GitHubRepository
		teams        []*models.GitHubTeam
		units        []*models.Unit
		reposFound   []string
		teamsFound   []string
		unitsFound   []string
	)
	// read the file and convert into standards list
	if err = structs.UnmarshalFile(path, &standards); err != nil {
		return
	}

	// truncate the standards table - as this is replaced each time
	err = crud.Truncate(ctx, adaptor, &models.GitHubRepositoryStandard{})
	if err != nil {
		return
	}
	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.All()...)
	if err != nil {
		return
	}

	// now get the unique repositories
	for _, std := range standards {
		var repo = (*models.GitHubRepository)(std.GitHubRepository)

		if !slices.Contains(reposFound, repo.UniqueValue()) {
			reposFound = append(reposFound, repo.UniqueValue())
			repositories = append(repositories, repo)
		}
	}

	// now grab the unique teams
	for _, repo := range repositories {
		var found = []*models.GitHubTeam{}
		// loop over each team
		for _, team := range repo.GitHubTeams {
			if !slices.Contains(teamsFound, team.UniqueValue()) {
				teamsFound = append(teamsFound, team.UniqueValue())
				teams = append(teams, team)
			}
			// now find the team in the existing set of teams
			for _, t := range teams {
				if t.UniqueValue() == team.UniqueValue() {
					found = append(found, t)
				}
			}
			// remap the teams so the ids are updated
			repo.GitHubTeams = found
		}

	}

	// get unique units
	for _, team := range teams {
		var found = []*models.Unit{}
		for _, unit := range team.Units {
			if !slices.Contains(unitsFound, unit.UniqueValue()) {
				unitsFound = append(unitsFound, unit.UniqueValue())
				units = append(units, unit)
			}
			// now find the unit in the existing set of
			for _, u := range units {
				if u.UniqueValue() == unit.UniqueValue() {
					found = append(found, u)
				}
			}
			team.Units = found
		}

	}

	// insert units
	if _, err = crud.Insert(ctx, adaptor, &models.Unit{}, units...); err != nil {
		return
	}
	// insert teams
	if _, err = crud.Insert(ctx, adaptor, &models.GitHubTeam{}, teams...); err != nil {
		return
	}
	// insert repositories
	if _, err = crud.Insert(ctx, adaptor, &models.GitHubRepository{}, repositories...); err != nil {
		return
	}

	// update the forien keys
	for _, std := range standards {
		std.GitHubRepositoryID = std.GitHubRepository.ID
	}
	// insert standards
	if _, err = crud.Insert(ctx, adaptor, &models.GitHubRepositoryStandard{}, standards...); err != nil {
		return
	}

	// joins
	// repo <-> team
	repoteams := []*models.GitHubRepositoryGitHubTeam{}
	for _, repo := range repositories {
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
		for _, unit := range team.Units {
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

	return

}
