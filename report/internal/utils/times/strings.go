package times

import "time"

func AsYMDString(t time.Time) string {
	return AsString(t, YMD)
}
func AsYMString(t time.Time) string {
	return AsString(t, YM)
}
func AsYString(t time.Time) string {
	return AsString(t, Y)
}

// AsString convert time string to a date formatted string
func AsString(t time.Time, layout Format) string {
	return t.Format(string(layout))
}
