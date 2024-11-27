package lib

import (
	"context"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/seed"
)

// processStandards handles importing github standards file with the structure of:
//   - GitHubRepositoryStandard
//     -- GitHubRepository
//     --- GitHubTeams
//     ---- Units
func processGithubStandards(ctx context.Context, adaptor dbs.Adaptor, path string) (res any, err error) {
	var standards []*models.GitHubRepositoryStandard
	// read the file and convert into standards list
	if err = structs.UnmarshalFile(path, &standards); err != nil {
		return
	}

	// truncate the standards table - as this is replaced each time
	err = crud.Truncate(ctx, adaptor, &models.GitHubRepositoryStandard{})
	if err != nil {
		return
	}
	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.All()...)
	if err != nil {
		return
	}

	res, err = seed.GitHubStandards(ctx, adaptor, standards)
	return

}
