package awscosts

import (
	"github.com/ministryofjustice/opg-reports/fakes"
)

// Fake generates a dummy version of the Cost struct with
// randomised values
// It uses 2023-12-01..2024-03-01 for the Date value (with RFC3339 format)
func Fake() *Cost {

	return &Cost{
		ID:           fakes.Int(100000, 999999),
		Ts:           fakes.DateAsStr(fakes.MinDate, fakes.MaxDate, fakes.DateFormat),
		Organisation: fakes.String(5),
		AccountID:    fakes.IntAsStr(100000000, 900000000),
		AccountName:  fakes.String(10),
		Unit:         fakes.String(12),
		Label:        fakes.String(12),
		Environment:  fakes.Choice[string]([]string{"production", "pre-production", "development"}),
		Service:      fakes.Choice[string]([]string{"ecs", "ec2", "r53", "s3", "sqs", "tax", "waf", "rds"}),
		Region:       fakes.Choice[string]([]string{"eu-west-1", "NoRegion", "eu-west-2", "us-east-1"}),
		Date:         fakes.DateAsStr(fakes.MinDate, fakes.MaxDate, fakes.DateFormat),
		Cost:         fakes.FloatAsStr(-5.0, 10),
	}

}

// FakeFrom is like Fake, but only adds fake data is it doesnt exist already on the cost
// struct passed
func FakeFrom(cost *Cost) *Cost {
	// generate a fake version of cost to start with
	faked := Fake()

	if cost.ID == 0 {
		cost.ID = faked.ID
	}
	if cost.Ts == "" {
		cost.Ts = faked.Ts
	}
	if cost.Organisation == "" {
		cost.Organisation = faked.Organisation
	}
	if cost.AccountID == "" {
		cost.AccountID = faked.AccountID
	}
	if cost.AccountName == "" {
		cost.AccountName = faked.AccountName
	}
	if cost.Unit == "" {
		cost.Unit = faked.Unit
	}
	if cost.Label == "" {
		cost.Label = faked.Label
	}
	if cost.Environment == "" {
		cost.Environment = faked.Environment
	}
	if cost.Service == "" {
		cost.Service = faked.Service
	}
	if cost.Region == "" {
		cost.Region = faked.Region
	}
	if cost.Date == "" {
		cost.Date = faked.Date
	}
	if cost.Cost == "" {
		cost.Cost = faked.Cost
	}

	return cost
}

// Fakes generates multiple fake Costs (based on count passed) as is a
// wrapper for Fake
func Fakes(count int) (fakes []*Cost) {
	fakes = []*Cost{}

	for i := 0; i < count; i++ {
		fakes = append(fakes, Fake())
	}
	return
}

func FakesFrom(count int, base *Cost) (fakes []*Cost) {
	fakes = []*Cost{}
	for i := 0; i < count; i++ {
		fakes = append(fakes, FakeFrom(base))
	}
	return
}
