package lib

import (
	"context"
	"fmt"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/seed"
)

// processAwsUptime handles importing github standards file with the structure of:
//   - AwsUptime
//     -- AwsAccount
//     -- Unit
func processAwsUptime(ctx context.Context, adaptor dbs.Adaptor, path string) (res any, err error) {
	var uptime []*models.AwsUptime

	if adaptor == nil {
		err = fmt.Errorf("adaptor is nil")
		return
	}
	// read the file and convert into standards list
	if err = structs.UnmarshalFile(path, &uptime); err != nil {
		return
	}

	res, err = seed.AwsUptime(ctx, adaptor, uptime)
	return

}
