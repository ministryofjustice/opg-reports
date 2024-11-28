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

func AwsAccounts(ctx context.Context, adaptor dbs.Adaptor, accounts []*models.AwsAccount) (res []*models.AwsAccount, err error) {
	var (
		units      []*models.Unit
		unitsFound []string
	)
	slog.Info("[seed] seeding aws accounts.")

	if adaptor == nil {
		err = fmt.Errorf("adaptor is nil")
		return
	}
	if !adaptor.Mode().Write() {
		err = fmt.Errorf("adaptor is not writable")
		return
	}

	// get the accounts
	for _, account := range accounts {
		if !slices.Contains(unitsFound, account.Unit.UniqueValue()) {
			unitsFound = append(unitsFound, account.Unit.UniqueValue())
			units = append(units, (*models.Unit)(account.Unit))
		}
	}

	// insert units
	if _, err = crud.Insert(ctx, adaptor, &models.Unit{}, units...); err != nil {
		return
	}

	// set the join
	for _, acc := range accounts {
		acc.UnitID = acc.Unit.ID
	}

	// insert accounts
	if _, err = crud.Insert(ctx, adaptor, &models.AwsAccount{}, accounts...); err != nil {
		return
	}

	res = accounts
	slog.Info("[seed] aws accounts done.")
	return
}
