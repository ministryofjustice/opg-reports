package reset

import (
	"opg-reports/report/packages/times"
	"testing"
	"time"
)

type scenario struct {
	Source   time.Time
	Interval times.Interval
	Expected time.Time
}

func TestPackagesResetTime(t *testing.T) {

	var tests = []scenario{
		{
			Interval: times.YEAR,
			Source:   times.FromString("2026-05-25T23:59:59"),
			Expected: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Interval: times.MONTH,
			Source:   times.FromString("2026-05-25T23:59:59"),
			Expected: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Interval: times.DAY,
			Source:   times.FromString("2026-02-28T23:59:59"),
			Expected: time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			Interval: times.HOUR,
			Source:   times.FromString("2026-02-28T23:59:59"),
			Expected: time.Date(2026, 2, 28, 23, 0, 0, 0, time.UTC),
		},
		{
			Interval: times.MINUTE,
			Source:   times.FromString("2026-02-28T23:59:59"),
			Expected: time.Date(2026, 2, 28, 23, 59, 0, 0, time.UTC),
		},
		{
			Interval: times.SECOND,
			Source:   times.FromString("2026-02-28T23:59:59"),
			Expected: time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC),
		},
	}

	for i, test := range tests {

		actual := Time(test.Source, test.Interval)
		if actual != test.Expected {
			t.Errorf("[%d] reset failed; expected [%v] actual [%v]", i, test.Expected, actual)
		}

	}

}
