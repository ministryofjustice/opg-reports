package awsu

import (
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

func Fake() (a *AwsUptime) {
	minDate, maxDate, dFormat := testhelpers.Dates()

	return &AwsUptime{
		ID:      fake.Int(100000, 999999),
		Ts:      fake.DateAsStr(minDate, maxDate, dFormat),
		Unit:    fake.String(12),
		Average: fake.Float(50.0, 100.0),
	}
}
