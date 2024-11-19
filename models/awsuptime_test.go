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
	_ dbs.Table           = &models.AwsUptime{}
	_ dbs.CreateableTable = &models.AwsUptime{}
	_ dbs.Insertable      = &models.AwsUptime{}
	_ dbs.Row             = &models.AwsUptime{}
	_ dbs.InsertableRow   = &models.AwsUptime{}
	_ dbs.Record          = &models.AwsUptime{}
)

// TestModelsAwsUptimeCRUD checks the github team table creation
// and inserting series of fake records works as expected
func TestModelsAwsUptimeCRUD(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err      error
		adaptor  *adaptors.Sqlite
		n        int                 = 100
		ctx      context.Context     = context.Background()
		dir      string              = t.TempDir()
		path     string              = filepath.Join(dir, "test.db")
		accounts []*models.AwsUptime = fakermany.Fake[*models.AwsUptime](n)
		results  []*models.AwsUptime
	)

	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	_, err = crud.CreateTable(ctx, adaptor, &models.AwsUptime{})
	if err != nil {
		t.Errorf("unexpected error for create table: [%s]", err.Error())
	}
	_, err = crud.CreateIndexes(ctx, adaptor, &models.AwsUptime{})
	if err != nil {
		t.Errorf("unexpected error for create indexes: [%s]", err.Error())
	}

	results, err = crud.Insert(ctx, adaptor, &models.AwsUptime{}, accounts...)
	if err != nil {
		t.Errorf("unexpected error for insert: [%s]", err.Error())
	}

	if len(results) != len(accounts) {
		t.Errorf("created records do not match expacted number - [%d] actual [%v]", len(accounts), len(results))
	}

}

var selectUptime = `
SELECT
	aws_uptime.*,
	json_object(
		'id', units.id,
		'name', units.name
	) as unit,
	 json_object(
		'id', aws_accounts.id,
		'number', aws_accounts.number,
		'name', aws_accounts.name,
		'label', aws_accounts.label,
		'environment', aws_accounts.environment
	) as aws_account
FROM aws_uptime
LEFT JOIN aws_accounts on aws_accounts.id = aws_uptime.aws_account_id
LEFT JOIN units on units.id = aws_accounts.unit_id
GROUP BY aws_uptime.id
ORDER BY aws_uptime.date ASC;
`

func TestModelsAwsUptimeUnitJoin(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err      error
		adaptor  *adaptors.Sqlite
		ctx      context.Context = context.Background()
		dir      string          = t.TempDir()
		path     string          = filepath.Join(dir, "test.db")
		uptime   []*models.AwsUptime
		accounts []*models.AwsAccount
		units    []*models.Unit
		results  []*models.AwsUptime
	)
	adaptor, err = adaptors.NewSqlite(path, false)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err.Error())
	}
	defer adaptor.DB().Close()

	// create test units
	units, err = testDBbuilder(ctx, adaptor, &models.Unit{}, fakermany.Fake[*models.Unit](3))
	if err != nil {
		t.Fatalf(err.Error())
	}
	// create test accounts
	accounts = fakermany.Fake[*models.AwsAccount](3)
	// now add a unit to each account
	for _, acc := range accounts {
		var unit = fakerextras.Choice(units)
		var unitjoin = models.UnitForeignKey(*unit)
		acc.UnitID = unit.ID
		acc.Unit = &unitjoin
	}
	// now save the accounts
	_, err = testDBbuilder(ctx, adaptor, &models.AwsAccount{}, accounts)
	if err != nil {
		t.Fatalf(err.Error())
	}
	uptime = fakermany.Fake[*models.AwsUptime](6)
	// now attach both account and unit to the uptimes
	for _, up := range uptime {
		var acc = fakerextras.Choice(accounts)
		var accjoin = models.AwsAccountForeignKey(*acc)
		up.AwsAccountID = acc.ID
		up.AwsAccount = &accjoin

	}
	// insert those
	_, err = testDBbuilder(ctx, adaptor, &models.AwsUptime{}, uptime)
	if err != nil {
		t.Fatalf(err.Error())
	}

	results, err = crud.Select[*models.AwsUptime](ctx, adaptor, selectUptime, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	//  check length matches
	if len(uptime) != len(results) {
		t.Errorf("length mismatch - expected [%d] actual [%v]", len(uptime), len(results))
	}

}
