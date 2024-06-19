package dates

import "time"

func Months(start time.Time, end time.Time) []time.Time {
	// reset the days to be the first day of each
	s := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
	e := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location())
	months := []time.Time{}

	for d := s; d.After(e) == false; d = d.AddDate(0, 1, 0) {
		months = append(months, d)
	}
	return months
}

func Days(start time.Time, end time.Time) []time.Time {
	days := []time.Time{}
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}
	return days
}
