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
	_ dbs.Table           = &models.Unit{}
	_ dbs.CreateableTable = &models.Unit{}
	_ dbs.Insertable      = &models.Unit{}
	_ dbs.Row             = &models.Unit{}
	_ dbs.InsertableRow   = &models.Unit{}
	_ dbs.Record          = &models.Unit{}
)

// TestModelsUnitCRUD checks the unit table and inserting series of fake
// records works as expected
func TestModelsUnitCRUD(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err     error
		adaptor *adaptors.Sqlite
		n       int             = 4
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
		units   []*models.Unit  = fakermany.Fake[*models.Unit](n)
		results []*models.Unit
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &models.Unit{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}
	_, err = crud.CreateIndexes(ctx, adaptor, &models.Unit{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &models.Unit{}, units...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(results) != len(units) {
		t.Errorf("created records do not match expacted number - [%d] actual [%v]", len(units), len(results))
	}

}

var selectUnitAndAccounts = `
SELECT
	units.*,
	json_group_array(
		json_object(
			'id', aws_accounts.id,
			'number', aws_accounts.number,
			'name', aws_accounts.name,
			'label', aws_accounts.label,
			'environment', aws_accounts.environment
		)
	) as aws_accounts
FROM units
LEFT JOIN aws_accounts ON aws_accounts.unit_id = units.id
GROUP BY units.id
ORDER BY units.name ASC;
`

func TestModelsUnitAwsAccountJoin(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err      error
		adaptor  *adaptors.Sqlite
		ctx      context.Context = context.Background()
		dir      string          = t.TempDir()
		path     string          = filepath.Join(dir, "test.db")
		accounts []*models.AwsAccount
		units    []*models.Unit
		results  []*models.Unit
	)
	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	defer adaptor.DB().Close()

	// create fake accounts & units
	accounts = fakermany.Fake[*models.AwsAccount](4)
	// generate the accounts
	units, err = testDBbuilder(ctx, adaptor, &models.Unit{}, fakermany.Fake[*models.Unit](3))
	if err != nil {
		t.Fatalf(err.Error())
	}

	// update the accounts to have unit id
	for _, account := range accounts {
		var unit = fakerextras.Choice(units)
		var u = models.UnitForeignKey(*unit)
		account.UnitID = unit.ID
		account.Unit = &u
	}
	// now save the accounts
	_, err = testDBbuilder(ctx, adaptor, &models.AwsAccount{}, accounts)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// fetch the data via select
	results, err = crud.Select[*models.Unit](ctx, adaptor, selectUnitAndAccounts, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(units) != len(results) {
		t.Errorf("incorrect number of records returned - expected [%d] actual [%v]", len(units), len(results))
	}
}
