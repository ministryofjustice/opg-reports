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

var defaultEnvironment string = "production"

// processAwsCosts handles importing github standards file with the structure of:
//   - AwsCost
//     -- AwsAccount
//     -- Unit
func processAwsCosts(ctx context.Context, adaptor dbs.Adaptor, path string) (res any, err error) {
	var (
		costs         []*models.AwsCost
		accounts      []*models.AwsAccount
		units         []*models.Unit
		accountsFound []string
		unitsFound    []string
	)
	if adaptor == nil {
		err = fmt.Errorf("adaptor is nil")
		return
	}

	// read the file and convert into standards list
	if err = structs.UnmarshalFile(path, &costs); err != nil {
		return
	}

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.All()...)
	if err != nil {
		return
	}

	// now get the unique repositories
	for _, cost := range costs {

		if !slices.Contains(accountsFound, cost.AwsAccount.UniqueValue()) {
			accountsFound = append(accountsFound, cost.AwsAccount.UniqueValue())
			accounts = append(accounts, (*models.AwsAccount)(cost.AwsAccount))
		}
		if !slices.Contains(unitsFound, cost.Unit.UniqueValue()) {
			unitsFound = append(accountsFound, cost.Unit.UniqueValue())
			units = append(units, (*models.Unit)(cost.Unit))
		}
	}

	// insert units
	if _, err = crud.Insert(ctx, adaptor, &models.Unit{}, units...); err != nil {
		return
	}
	// update the unit id on the account
	for _, acc := range accounts {
		var found *models.Unit
		for _, u := range units {
			if acc.Unit.UniqueValue() == u.UniqueValue() {
				found = u
			}
		}
		if found != nil {
			acc.UnitID = found.ID
			acc.Unit = (*models.UnitForeignKey)(found)
		}
		// fix the environment value
		if acc.Environment == "" || acc.Environment == "null" {
			acc.Environment = defaultEnvironment
		}
	}
	// insert accounts
	if _, err = crud.Insert(ctx, adaptor, &models.AwsAccount{}, accounts...); err != nil {
		return
	}

	// pretty.Print(accounts)

	// now update the account id join field on the uptime record
	for _, cost := range costs {
		for _, acc := range accounts {
			if cost.AwsAccount.UniqueValue() == acc.UniqueValue() {
				cost.AwsAccountID = acc.ID
			}
		}
	}
	// insert uptime
	if _, err = crud.Insert(ctx, adaptor, &models.AwsCost{}, costs...); err != nil {
		return
	}
	res = costs
	return

}
