package lib

import (
	"context"
	"fmt"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/seed"
)

// processAwsCosts handles importing github standards file with the structure of:
//   - AwsCost
//     -- AwsAccount
//     -- Unit
func processAwsCosts(ctx context.Context, adaptor dbs.Adaptor, path string) (res any, err error) {
	var costs []*models.AwsCost

	if adaptor == nil {
		err = fmt.Errorf("adaptor is nil")
		return
	}

	// read the file and convert into standards list
	if err = structs.UnmarshalFile(path, &costs); err != nil {
		return
	}

	res, err = seed.AwsCosts(ctx, adaptor, costs)
	return

}
