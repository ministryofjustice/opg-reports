package awsc

import (
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

func Fake() (a *AwsCost) {
	minDate, maxDate, dFormat := testhelpers.Dates()

	return &AwsCost{
		ID:           fake.Int(100000, 999999),
		Ts:           fake.DateAsStr(minDate, maxDate, dFormat),
		Organisation: fake.String(5),
		AccountID:    fake.IntAsStr(100000000, 900000000),
		AccountName:  fake.String(10),
		Unit:         fake.String(12),
		Label:        fake.String(12),
		Environment:  fake.Choice[string]([]string{"production", "pre-production", "development"}),
		Service:      fake.Choice[string]([]string{"ecs", "ec2", "r53", "s3", "sqs", "tax", "waf", "rds"}),
		Region:       fake.Choice[string]([]string{"eu-west-1", "NoRegion", "eu-west-2", "us-east-1"}),
		Date:         fake.DateAsStr(minDate, maxDate, dFormat),
		Cost:         fake.FloatAsStr(-500.0, 15000),
	}

}
