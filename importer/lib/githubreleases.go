package lib

import (
	"context"
	"fmt"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/seed"
)

// processGithubReleases handles importing github standards file with the data structure of:
//   - GitHubRelease
//     -- GitHubRepository
//     ---- GitHubTeam
func processGithubReleases(ctx context.Context, adaptor dbs.Adaptor, path string) (res any, err error) {
	var releases []*models.GitHubRelease
	if adaptor == nil {
		err = fmt.Errorf("adaptor is nil")
		return
	}
	// read the file and convert into standards list
	if err = structs.UnmarshalFile(path, &releases); err != nil {
		return
	}

	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.All()...)
	if err != nil {
		return
	}

	res, err = seed.GitHubReleases(ctx, adaptor, releases)
	res = releases
	return

}
