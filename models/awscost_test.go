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
	_ dbs.Table           = &models.AwsCost{}
	_ dbs.CreateableTable = &models.AwsCost{}
	_ dbs.Insertable      = &models.AwsCost{}
	_ dbs.Row             = &models.AwsCost{}
	_ dbs.InsertableRow   = &models.AwsCost{}
	_ dbs.Record          = &models.AwsCost{}
)

// TestModelsAwsCostCRUD checks the unit table and inserting series of fake
// records works as expected
func TestModelsAwsCostCRUD(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err     error
		adaptor *adaptors.Sqlite
		n       int               = 4
		ctx     context.Context   = context.Background()
		dir     string            = t.TempDir()
		path    string            = filepath.Join(dir, "test.db")
		units   []*models.AwsCost = fakermany.Fake[*models.AwsCost](n)
		results []*models.AwsCost
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &models.AwsCost{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}
	_, err = crud.CreateIndexes(ctx, adaptor, &models.AwsCost{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &models.AwsCost{}, units...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(results) != len(units) {
		t.Errorf("created records do not match expacted number - [%d] actual [%v]", len(units), len(results))
	}

}

var selectCostsAndAccounts = `
SELECT
	aws_costs.*,
	json_object(
		'id', units.id,
		'name', units.name
	) as unit,
	json_object(
		'id', aws_accounts.id,
		'number', aws_accounts.number,
		'name', aws_accounts.name,
		'label', aws_accounts.label,
		'environment', aws_accounts.environment,
		'unit_id', aws_accounts.unit_id
	) as aws_account
FROM aws_costs
LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
GROUP BY aws_costs.id
ORDER BY aws_costs.date ASC;
`

func TestModelsAwsCostJoins(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err      error
		adaptor  *adaptors.Sqlite
		ctx      context.Context = context.Background()
		dir      string          = t.TempDir()
		path     string          = filepath.Join(dir, "test.db")
		costs    []*models.AwsCost
		accounts []*models.AwsAccount
		units    []*models.Unit
		results  []*models.AwsCost
	)
	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	defer adaptor.DB().Close()

	// create fake units and save to db
	units = fakermany.Fake[*models.Unit](3)
	_, err = testDBbuilder(ctx, adaptor, &models.Unit{}, units)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// create the accounts and attach the units to them before saving
	accounts = fakermany.Fake[*models.AwsAccount](4)
	for _, account := range accounts {
		var unit = fakerextras.Choice(units)
		var u = models.UnitForeignKey(*unit)
		account.UnitID = unit.ID
		account.Unit = &u
	}
	_, err = testDBbuilder(ctx, adaptor, &models.AwsAccount{}, accounts)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// now create series of dummy costs and attach the account
	costs = fakermany.Fake[*models.AwsCost](15)

	for _, cost := range costs {
		var acc = fakerextras.Choice(accounts)
		var a = models.AwsAccountForeignKey(*acc)

		cost.AwsAccountID = acc.ID
		cost.AwsAccount = &a
	}

	// insert the costs
	_, err = testDBbuilder(ctx, adaptor, &models.AwsCost{}, costs)
	if err != nil {
		t.Fatalf(err.Error())
	}

	results, err = crud.Select[*models.AwsCost](ctx, adaptor, selectCostsAndAccounts, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// check length
	if len(costs) != len(results) {
		t.Errorf("incorrect number of records returned - expected [%d] actual [%v]", len(units), len(results))
	}

	// check content
	for _, item := range results {
		if item.AwsAccountID != item.AwsAccount.ID {
			t.Errorf("join failed on accounts")
		}
		if item.Unit.ID != item.AwsAccount.UnitID {
			t.Errorf("join failed on account -> unit ")
		}
	}
}
