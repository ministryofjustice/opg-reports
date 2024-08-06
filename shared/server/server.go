// Package server provides a series of interfaces and concrete structs for request handling
package server

import (
	"opg-reports/shared/dates"
	"time"
)

func GetStartEndDates(parameters map[string][]string) (time.Time, time.Time) {
	now := time.Now().UTC()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startDate := firstDay
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if starts, ok := parameters["start"]; ok && starts[0] != "-" {
		startDate, _ = dates.StringToDateDefault(starts[0], "-", firstDay.Format(dates.FormatYM))
	}
	if ends, ok := parameters["end"]; ok && ends[0] != "-" {
		endDate, _ = dates.StringToDateDefault(ends[0], "-", endDate.Format(dates.FormatYM))
	}
	return startDate, endDate
}
