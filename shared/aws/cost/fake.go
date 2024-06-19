package cost

import (
	"opg-reports/shared/fake"
	"time"
)

// Fake returns a generated Cost item using fake data
func Fake(c *Cost, minDate time.Time, maxDate time.Time, dFormat string) (f *Cost) {
	if c == nil {
		c = New(nil)
	}
	if c.AccountEnvironment == "" {
		c.AccountEnvironment = fake.String(3)
	}
	if c.AccountId == "" {
		c.AccountId = fake.IntAsStr(100000, 999999)
	}
	if c.AccountLabel == "" {
		c.AccountLabel = fake.String(5)
	}
	if c.AccountName == "" {
		c.AccountName = fake.String(10)
	}
	if c.AccountOrganisation == "" {
		c.AccountOrganisation = fake.String(3)
	}
	if c.AccountUnit == "" {
		c.AccountUnit = fake.String(4)
	}
	if c.Service == "" {
		c.Service = fake.String(8)
	}
	if c.Region == "" {
		c.Region = fake.String(4)
	}
	if c.Cost == "" {
		c.Cost = fake.FloatAsStr(1, 10000)
	}
	if c.Date == "" {
		c.Date = fake.DateAsStr(minDate, maxDate, dFormat)
	}

	f = c
	return
}
