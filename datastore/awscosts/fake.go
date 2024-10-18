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
		Cost:         fakes.FloatAsStr(-500.0, 15000),
	}

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
