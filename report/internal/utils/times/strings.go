package times

import "time"

// AsString convert time string to a date formatted string
func AsString(t time.Time, layout Format) string {
	return t.Format(string(layout))
}
