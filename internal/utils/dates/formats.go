package dates

import "time"

// GuessDateFormat will return the format to use for the date string passed, using
// time.RFC3339 as base. used in time.Format calls
//
// Passing 2024 would return 2006, passing 2024-12-01 would return 2006-01-02
// and so on
func GuessDateFormat(value string) string {
	var layout = time.RFC3339
	max := len(layout)
	l := len(value)
	if l > max {
		return layout
	}
	f := layout[:l]
	return f
}
