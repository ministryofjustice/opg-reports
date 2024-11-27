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

// AwsUptime handles importing github standards file with the structure of:
//   - AwsUptime
//     -- AwsAccount
//     -- Unit
func AwsUptime(ctx context.Context, adaptor dbs.Adaptor, uptime []*models.AwsUptime) (res []*models.AwsUptime, err error) {
	var (
		accounts           []*models.AwsAccount
		units              []*models.Unit
		accountsFound      []string
		unitsFound         []string
		defaultEnvironment string = "production"
	)
	slog.Info("[seed] seeding aws uptime.")
	if adaptor == nil {
		err = fmt.Errorf("adaptor is nil")
		return
	}
	if !adaptor.Mode().Write() {
		err = fmt.Errorf("adaptor is not writable")
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
		// fix the environment value
		if account.Environment == "" || account.Environment == "null" {
			account.Environment = defaultEnvironment
		}
	}

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
	slog.Info("[seed] aws uptime done.")
	return

}
