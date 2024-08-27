package datastore

import (
	"github.com/ministryofjustice/opg-reports/datastore/aws_costs/awsc"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
)

type Record interface {
	awsc.AwsCost | awsc.AwsCostsTracker | ghs.GithubStandard | ghs.GithubStandardsTracker
}
