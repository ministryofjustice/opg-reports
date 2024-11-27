package seed

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/models"
)

// GitHubTeams handles importing with structure of:
//   - GitHubTeams
//     -- GitHubRepositories
//     -- Units
func GitHubTeams(ctx context.Context, adaptor dbs.Adaptor, teams []*models.GitHubTeam) (res []*models.GitHubTeam, err error) {
	var (
		units      []*models.Unit
		repos      []*models.GitHubRepository
		reposFound []string
		unitsFound []string
	)
	slog.Info("[seed] seeding github teams.")

	if adaptor == nil {
		err = fmt.Errorf("adaptor is nil")
		return
	}
	if !adaptor.Mode().Write() {
		err = fmt.Errorf("adaptor is not writable")
		return
	}
	// now get the unique repos and units
	for _, team := range teams {
		// find the unique repositories
		for _, repo := range team.GitHubRepositories {
			if !slices.Contains(reposFound, repo.UniqueValue()) {
				reposFound = append(reposFound, repo.UniqueValue())
				repos = append(repos, repo)
			}
		}
		// find the unit units
		for _, unit := range team.Units {
			if !slices.Contains(unitsFound, unit.UniqueValue()) {
				unitsFound = append(unitsFound, unit.UniqueValue())
				units = append(units, unit)
			}
		}
	}
	// now re map the team->repo & team->unit setup on the models
	for _, team := range teams {
		var (
			teamRepos = []*models.GitHubRepository{}
			teamUnits = []*models.Unit{}
		)
		// get the repos for this team from the found set (pointers)
		for _, repo := range team.GitHubRepositories {
			for _, r := range repos {
				if repo.UniqueValue() == r.UniqueValue() {
					teamRepos = append(teamRepos, r)
				}
			}
		}
		// get the units
		for _, unit := range team.Units {
			for _, u := range units {
				if unit.UniqueValue() == u.UniqueValue() {
					teamUnits = append(teamUnits, u)
				}
			}
		}
		team.GitHubRepositories = teamRepos
		team.Units = teamUnits

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
	if _, err = crud.Insert(ctx, adaptor, &models.GitHubRepository{}, repos...); err != nil {
		return
	}

	// joins
	// repo <-> team
	repoteams := []*models.GitHubRepositoryGitHubTeam{}
	for _, team := range teams {
		for _, repo := range team.GitHubRepositories {
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

	res = teams
	// pretty.Print(res)
	slog.Info("[seed] github teams done.")
	return
}
