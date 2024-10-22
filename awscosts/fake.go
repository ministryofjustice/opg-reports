package awscosts

import "github.com/ministryofjustice/opg-reports/fakes"

func Fake(base *Cost) (cost *Cost) {
	cost = &Cost{}

	if base != nil {
		*cost = *base
	}

	if cost.ID == 0 {
		cost.ID = fakes.Int(100000, 999999)
	}
	if cost.Ts == "" {
		cost.Ts = fakes.DateAsStr(fakes.MinDate, fakes.MaxDate, fakes.DateFormat)
	}
	if cost.Organisation == "" {
		cost.Organisation = fakes.String(5)
	}
	if cost.AccountID == "" {
		cost.AccountID = fakes.IntAsStr(100000000, 900000000)
	}
	if cost.AccountName == "" {
		cost.AccountName = fakes.String(10)
	}
	if cost.Unit == "" {
		cost.Unit = fakes.String(12)
	}
	if cost.Label == "" {
		cost.Label = fakes.String(12)
	}
	if cost.Environment == "" {
		cost.Environment = fakes.Choice([]string{"production", "pre-production", "development"})
	}
	if cost.Service == "" {
		cost.Service = fakes.Choice([]string{"ecs", "ec2", "r53", "s3", "sqs", "tax", "waf", "rds"})
	}
	if cost.Region == "" {
		cost.Region = fakes.Choice([]string{"eu-west-1", "NoRegion", "eu-west-2", "us-east-1"})
	}
	if cost.Date == "" {
		cost.Date = fakes.DateAsStr(fakes.MinDate, fakes.MaxDate, fakes.DateFormat)
	}
	if cost.Cost == "" {
		cost.Cost = fakes.FloatAsStr(-4.5, 10.5)
	}

	return
}

func Fakes(count int, cost *Cost) (fakes []*Cost) {
	fakes = []*Cost{}

	for i := 0; i < count; i++ {
		f := Fake(cost)
		fakes = append(fakes, f)
	}
	return
}
