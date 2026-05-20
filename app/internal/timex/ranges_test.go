package timex

import (
	"testing"
	"time"
)

type rangeTest struct {
	Start    time.Time
	End      time.Time
	Interval Interval
	Inc      int
	Expected []time.Time
}

func TestTimexRanges(t *testing.T) {

	var tests = []*rangeTest{
		// DAYS
		// 	2024/02/27 - 2024/03/1
		// 	leap year
		{
			Start:    time.Date(2024, 2, 27, 23, 55, 59, 59, time.UTC),
			End:      time.Date(2024, 3, 1, 1, 0, 0, 1, time.UTC),
			Interval: DAY,
			Inc:      1,
			Expected: []time.Time{
				time.Date(2024, 2, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		// 	2024/02/27 - 2024/03/1
		// 	leap year, skipping days
		{
			Start:    time.Date(2024, 2, 27, 23, 55, 59, 59, time.UTC),
			End:      time.Date(2024, 3, 1, 1, 0, 0, 2, time.UTC),
			Interval: DAY,
			Inc:      5,
			Expected: []time.Time{
				time.Date(2024, 2, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for i, test := range tests {
		actuals := Range(test.Start, test.End, test.Interval, test.Inc)
		// now loop over actuals and compare
		for j, actual := range actuals {
			var expected = test.Expected[j]

			if !actual.Equal(expected) {
				t.Errorf("[test:%d][%d] actual did nto match expected - [%v] == [%v]", i, j, actual, expected)
			}

		}
	}

}
