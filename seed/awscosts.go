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

// AwsCosts handles importing github standards file with the structure of:
//   - AwsCost
//     -- AwsAccount
//     -- Unit
func AwsCosts(ctx context.Context, adaptor dbs.Adaptor, costs []*models.AwsCost) (res []*models.AwsCost, err error) {
	var (
		accounts           []*models.AwsAccount
		units              []*models.Unit
		accountsFound      []string
		unitsFound         []string
		defaultEnvironment string = "production"
	)

	slog.Info("[seed] seeding aws costs.")
	if adaptor == nil {
		err = fmt.Errorf("adaptor is nil")
		return
	}
	if !adaptor.Mode().Write() {
		err = fmt.Errorf("adaptor is not writable")
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
	slog.Info("[seed] aws costs done.")
	return
}
