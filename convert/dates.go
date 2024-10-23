package convert

import "time"

func DateAddYear(d time.Time, mod int) time.Time {
	return d.AddDate(mod, 0, 0)
}

func DateAddMonth(d time.Time, mod int) time.Time {
	return d.AddDate(0, mod, 0)
}
func DateAddDay(d time.Time, mod int) time.Time {
	return d.AddDate(0, 0, mod)
}

// ResetMonth resets the day of the month to the 1st and zeros the time
func DateResetYear(d time.Time) time.Time {
	return time.Date(d.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
}

// ResetMonth resets the day of the month to the 1st and zeros the time
func DateResetMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// ResetDay resets time of the day to zeros
func DateResetDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}
