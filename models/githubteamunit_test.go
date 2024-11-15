package models_test

import (
	"context"
	"fmt"
	"path/filepath"
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
	_ dbs.Table           = &models.GitHubTeamUnit{}
	_ dbs.CreateableTable = &models.GitHubTeamUnit{}
	_ dbs.Insertable      = &models.GitHubTeamUnit{}
	_ dbs.Row             = &models.GitHubTeamUnit{}
	_ dbs.InsertableRow   = &models.GitHubTeamUnit{}
	_ dbs.Record          = &models.GitHubTeamUnit{}
)

func testDBbuilder[T dbs.TableOfRecord](ctx context.Context, adaptor *adaptors.Sqlite, itemType T, insert []T) (results []T, err error) {

	_, err = crud.CreateTable(ctx, adaptor, itemType)
	if err != nil {
		return
	}

	_, err = crud.CreateIndexes(ctx, adaptor, itemType)
	if err != nil {
		return
	}

	results, err = crud.Insert(ctx, adaptor, itemType, insert...)
	if err != nil {
		return
	}

	if len(results) != len(insert) {
		err = fmt.Errorf("inserted record count mistmatch")
	}
	return
}

var selectUnits string = `
SELECT
	units.*,
	json_group_array(json_object('id', github_teams.id,'name', github_teams.name)) as github_teams
FROM units
LEFT JOIN github_teams_units on github_teams_units.unit_id = units.id
LEFT JOIN github_teams on github_teams.id = github_teams_units.github_team_id
GROUP BY units.id
ORDER BY units.name ASC;
`

// TestModelsGithubTeamUnitJoin checks the join logic from
// unit->githubteams is working
func TestModelsGithubTeamUnitJoin(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err         error
		adaptor     *adaptors.Sqlite
		ctx         context.Context = context.Background()
		dir         string          = t.TempDir()
		path        string          = filepath.Join(dir, "test.db")
		units       []*models.Unit
		resultUnits []*models.Unit
		teams       []*models.GitHubTeam
		joins       []*models.GitHubTeamUnit
	)
	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	defer adaptor.DB().Close()

	units, err = testDBbuilder(ctx, adaptor, &models.Unit{}, fakermany.Fake[*models.Unit](4))
	if err != nil {
		t.Fatalf(err.Error())
	}
	teams, err = testDBbuilder(ctx, adaptor, &models.GitHubTeam{}, fakermany.Fake[*models.GitHubTeam](6))
	if err != nil {
		t.Fatalf(err.Error())
	}

	// now we create the joins and insert them
	for _, unit := range units {
		var list = fakerextras.Choices(teams, 1)
		// set the team on the unit
		unit.GitHubTeams = list
		for _, gt := range list {
			joins = append(joins, &models.GitHubTeamUnit{UnitID: unit.ID, GithubTeamID: gt.ID})
		}
	}
	// insert the joins and create the table etc
	_, err = testDBbuilder(ctx, adaptor, &models.GitHubTeamUnit{}, joins)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// now select the units to and see if the teams are included!
	resultUnits, err = crud.Select[*models.Unit](ctx, adaptor, selectUnits, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// now check the results contain the correct join data
	for _, result := range resultUnits {
		// grab actually returned teams
		var actualTeams = result.GitHubTeams
		// get the teams that were generated for this unit
		var expectedTeams = []*models.GitHubTeam{}
		for _, unit := range units {
			if unit.ID == result.ID {
				expectedTeams = unit.GitHubTeams
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
