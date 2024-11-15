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
	"github.com/ministryofjustice/opg-reports/internal/pretty"
	"github.com/ministryofjustice/opg-reports/models"
)

// Interface checks
var (
	_ dbs.Table           = &models.Organisation{}
	_ dbs.CreateableTable = &models.Organisation{}
	_ dbs.Insertable      = &models.Organisation{}
	_ dbs.Row             = &models.Organisation{}
	_ dbs.InsertableRow   = &models.Organisation{}
	_ dbs.Record          = &models.Organisation{}
)

var selectOrgUnits = `
SELECT
	organisations.*,
	json_group_array(json_object('id', units.id,'name', units.name)) as units
FROM organisations
LEFT JOIN units on units.organisation_id = organisations.id
GROUP BY units.id
ORDER BY units.name ASC;
`

func TestModelsOrganisationUnitJoin(t *testing.T) {
	fakerextras.AddProviders()
	var (
		n       int = 2
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
		orgs    []*models.Organisation
		// units          []*models.Unit
		generatedUnits []*models.Unit
		results        []*models.Organisation
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}

	orgs, err = testDBbuilder(ctx, adaptor, &models.Organisation{}, fakermany.Fake[*models.Organisation](n))
	if err != nil {
		t.Fatalf(err.Error())
	}

	generatedUnits = fakermany.Fake[*models.Unit](n)
	// now we create the joins and insert them
	for i, org := range orgs {
		// set the team on the unit
		org.Units = []*models.Unit{generatedUnits[i]}
		for _, unit := range org.Units {
			unit.OrganisationID = org.ID
		}
	}

	_, err = testDBbuilder(ctx, adaptor, &models.Unit{}, generatedUnits)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// now select the units to and see if the teams are included!
	results, err = crud.Select[*models.Organisation](ctx, adaptor, selectOrgUnits, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	pretty.Print(results)
}

// TestModelsOrganisationCRUD checks the unit table and inserting series of fake
// records works as expected
func TestModelsOrganisationCRUD(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err     error
		adaptor *adaptors.Sqlite
		n       int                    = 2
		ctx     context.Context        = context.Background()
		dir     string                 = t.TempDir()
		path    string                 = filepath.Join(dir, "test.db")
		units   []*models.Organisation = fakermany.Fake[*models.Organisation](n)
		results []*models.Organisation
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &models.Organisation{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}
	_, err = crud.CreateIndexes(ctx, adaptor, &models.Organisation{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &models.Organisation{}, units...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(results) != len(units) {
		t.Errorf("created records do not match expacted number - [%d] actual [%v]", len(units), len(results))
	}

}
