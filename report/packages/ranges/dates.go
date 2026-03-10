package ranges

import (
	"opg-reports/report/packages/reset"
	"opg-reports/report/packages/times"
	"opg-reports/report/packages/types/interfaces"
	"time"
)

// Months will create a list of months between the start and
// end date provided of type T
func Months[T interfaces.DateTypes](start T, end T) (set []T) {
	var (
		is     T
		months []time.Time = []time.Time{}
		s      time.Time   = times.Time(start)
		e      time.Time   = times.Time(end)
	)
	set = []T{}
	// reset the start date to the first of the month
	s = reset.Month(&s)
	// add 1 month to the end date reset the month and remove 1 second to get the
	// last second of the month
	e = e.AddDate(0, 1, 0)
	e = reset.Month(&e).Add(-(time.Second * 1))
	// create the raw month values
	for d := s; d.After(e) == false; d = d.AddDate(0, 1, 0) {
		months = append(months, d)
	}
	// setup the return values
	switch any(is).(type) {
	case time.Time:
		for _, m := range months {
			var i interface{} = m
			set = append(set, i.(T))
		}
	case string:
		for _, m := range months {
			var i interface{} = m.Format(times.YM)
			set = append(set, i.(T))
		}
	}

	return
}
