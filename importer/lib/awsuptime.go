package lib

import (
	"context"
	"slices"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
)

// processAwsUptime handles importing github standards file with the structure of:
//   - AwsUptime
//     -- AwsAccount
//     -- Unit
func processAwsUptime(ctx context.Context, adaptor dbs.Adaptor, path string) (res any, err error) {
	var (
		uptime        []*models.AwsUptime
		accounts      []*models.AwsAccount
		units         []*models.Unit
		accountsFound []string
		unitsFound    []string
	)
	// read the file and convert into standards list
	if err = structs.UnmarshalFile(path, &uptime); err != nil {
		return
	}

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.All()...)
	if err != nil {
		return
	}

	// now get the unique repositories
	for _, up := range uptime {
		var (
			account *models.AwsAccount = (*models.AwsAccount)(up.AwsAccount)
			unit    *models.Unit       = (*models.Unit)(up.Unit)
		)
		if !slices.Contains(accountsFound, account.UniqueValue()) {
			accountsFound = append(accountsFound, account.UniqueValue())
			accounts = append(accounts, account)
		}
		if !slices.Contains(unitsFound, unit.UniqueValue()) {
			unitsFound = append(unitsFound, unit.UniqueValue())
			units = append(units, unit)
		}
		// now add the unit to the account
		account.Unit = (*models.UnitForeignKey)(unit)
	}

	// pretty.Print(uptime)

	// insert units
	if _, err = crud.Insert(ctx, adaptor, &models.Unit{}, units...); err != nil {
		return
	}
	// update the id on the account model with the unit id
	for _, acc := range accounts {
		acc.UnitID = acc.Unit.ID
	}
	// insert accounts
	if _, err = crud.Insert(ctx, adaptor, &models.AwsAccount{}, accounts...); err != nil {
		return
	}
	// now update the account id join field on the uptime record
	for _, up := range uptime {
		for _, acc := range accounts {
			if up.AwsAccount.UniqueValue() == acc.UniqueValue() {
				up.AwsAccountID = acc.ID
			}
		}
	}
	// insert uptime
	if _, err = crud.Insert(ctx, adaptor, &models.AwsUptime{}, uptime...); err != nil {
		return
	}
	res = uptime
	return

}
