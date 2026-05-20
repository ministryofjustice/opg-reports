package timex

import "time"

// GetFormat will return the format to use for the date string passed, using
// time.RFC3339 as base.
//
// Passing 2024 would return 2006, passing 2024-12-01 would return 2006-01-02
// and so on
//
// Used in coverting between string and time.Time for the Parse call
func GetFormat(value string) string {
	var (
		layout string = string(FULL)
		max    int    = len(layout)
		l      int    = len(value)
	)
	// return FULL version f the length is longer
	if l > max {
		return layout
	}
	return layout[:l]
}

// ToString converts time to a string using a fixed layout
func ToString(t time.Time, layout Format) (s string) {
	s = t.Format(string(layout))
	return
}
