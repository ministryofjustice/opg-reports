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
	_ dbs.Table           = &models.AwsAccount{}
	_ dbs.CreateableTable = &models.AwsAccount{}
	_ dbs.Insertable      = &models.AwsAccount{}
	_ dbs.Row             = &models.AwsAccount{}
	_ dbs.InsertableRow   = &models.AwsAccount{}
	_ dbs.Record          = &models.AwsAccount{}
)

// TestModelsAwsAccountCRUD checks the github team table creation
// and inserting series of fake records works as expected
func TestModelsAwsAccountCRUD(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err      error
		adaptor  *adaptors.Sqlite
		n        int                  = 100
		ctx      context.Context      = context.Background()
		dir      string               = t.TempDir()
		path     string               = filepath.Join(dir, "test.db")
		accounts []*models.AwsAccount = fakermany.Fake[*models.AwsAccount](n)
		results  []*models.AwsAccount
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &models.AwsAccount{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}
	_, err = crud.CreateIndexes(ctx, adaptor, &models.AwsAccount{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &models.AwsAccount{}, accounts...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(results) != len(accounts) {
		t.Errorf("created records do not match expacted number - [%d] actual [%v]", len(accounts), len(results))
	}

}

var selectAccounts = `
SELECT
	aws_accounts.*,
	json_object(
		'id', units.id,
		'name', units.name
	) as unit
FROM aws_accounts
LEFT JOIN units on units.id = aws_accounts.unit_id
GROUP BY aws_accounts.id
ORDER BY aws_accounts.name ASC;
`

func TestModelsAwsAccountUnitJoin(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err      error
		adaptor  *adaptors.Sqlite
		ctx      context.Context = context.Background()
		dir      string          = t.TempDir()
		path     string          = filepath.Join(dir, "test.db")
		accounts []*models.AwsAccount
		units    []*models.Unit
		results  []*models.AwsAccount
	)
	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	defer adaptor.DB().Close()

	// create test repos
	units, err = testDBbuilder(ctx, adaptor, &models.Unit{}, fakermany.Fake[*models.Unit](3))
	if err != nil {
		t.Fatalf(err.Error())
	}

	// create test releaases
	accounts = fakermany.Fake[*models.AwsAccount](5)
	// now add a repo to each release
	for _, acc := range accounts {
		var unit = fakerextras.Choice(units)
		var join = models.UnitForeignKey(*unit)
		acc.UnitID = unit.ID
		acc.Unit = &join
	}
	// now save the items
	_, err = testDBbuilder(ctx, adaptor, &models.AwsAccount{}, accounts)
	if err != nil {
		t.Fatalf(err.Error())
	}

	results, err = crud.Select[*models.AwsAccount](ctx, adaptor, selectAccounts, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	//  check length matches
	if len(accounts) != len(results) {
		t.Errorf("length mismatch - expected [%d] actual [%v]", len(accounts), len(results))
	}

	// now check the selected results match the generated version
	for _, res := range results {
		var release *models.AwsAccount
		for _, r := range accounts {
			if r.ID == res.ID {
				release = r
			}
		}
		if release == nil {
			t.Errorf("failed to find release")
		}
		if release.UnitID != res.UnitID {
			t.Errorf("repo ID mismatch")
		}
		if release.Unit.ID != res.Unit.ID {
			t.Errorf("repo ID mismatch")
		}
	}

}
