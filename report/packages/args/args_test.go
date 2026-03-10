package args

import (
	"opg-reports/report/packages/times"
	"testing"
	"time"
)

type tDateArgs struct {
	Source   time.Time
	Expected *Dates
}

func TestPackagesArgsDates(t *testing.T) {
	// var now =
	var tests = []*tDateArgs{
		// test logic so on the first day of a month the default dates
		// start should be the first day of the previous month to ensure
		// all data for the previous month is included
		{
			Source: times.FromString("2025-11-01T01:59:50"),
			Expected: &Dates{
				End:        times.FromString("2025-11-01T00:00:00"),
				Start:      times.FromString("2025-10-01T00:00:00"),
				StartCosts: times.FromString("2025-08-01T00:00:00"),
			},
		},
		// early march date, testing jan reset
		{
			Source: times.FromString("2026-03-05T01:59:50"),
			Expected: &Dates{
				End:        times.FromString("2026-03-05T00:00:00"),
				Start:      times.FromString("2026-03-01T00:00:00"),
				StartCosts: times.FromString("2026-01-01T00:00:00"),
			},
		},
		// feb date to test the year overlap for costs
		{
			Source: times.FromString("2026-02-20T05:59:50"),
			Expected: &Dates{
				End:        times.FromString("2026-02-20T00:00:00"),
				Start:      times.FromString("2026-02-01T00:00:00"),
				StartCosts: times.FromString("2025-12-01T00:00:00"),
			},
		},
	}

	for i, test := range tests {
		actual := Default[*Dates](test.Source)
		if actual.End != test.Expected.End {
			t.Errorf("[%d] END dates dont match, expected [%v] actual [%v]", i, test.Expected.End, actual.End)
		}
		if actual.Start != test.Expected.Start {
			t.Errorf("[%d] START dates dont match, expected [%v] actual [%v]", i, test.Expected.Start, actual.Start)
		}
		if actual.StartCosts != test.Expected.StartCosts {
			t.Errorf("[%d] COST START dates dont match, expected [%v] actual [%v]", i, test.Expected.StartCosts, actual.StartCosts)
		}
	}

}
