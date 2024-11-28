// dateintervals is a small package that contains all the
// commonly used date intervals within the project
// represented as constants
package dateintervals

import "github.com/ministryofjustice/opg-reports/internal/dateformats"

type Interval string

const (
	Year  Interval = "year"
	Month Interval = "month"
	Day   Interval = "day"
)

func Format(interval Interval) (format string) {
	switch interval {
	case Year:
		format = dateformats.Y
	case Month:
		format = dateformats.YM
	case Day:
		format = dateformats.YMD
	}
	return
}
