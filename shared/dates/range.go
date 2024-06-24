package dates

import (
	"slices"
	"time"
)

// Months returns a slice of time.Time for every month between the start & end values passed
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

// InMonth returns true if the YYYY-MM date string passed is within the list of months
func InMonth(dateStr string, months []string) bool {
	if d, e := StringToDate(dateStr); e == nil && slices.Contains(months, d.Format(FormatYM)) {
		return true
	}
	return false
}

// Days returns a slice of all days within the start and end dates
func Days(start time.Time, end time.Time) []time.Time {
	days := []time.Time{}
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}
	return days
}
