package times

import (
	"time"
)

func AsYMDString(t time.Time) string {
	return AsString(t, YMD)
}
func AsYMString(t time.Time) string {
	return AsString(t, YM)
}
func AsYString(t time.Time) string {
	return AsString(t, Y)
}

func AsYMDStrings(list []time.Time) []string {
	return AsStrings(list, YMD)
}
func AsYMStrings(list []time.Time) []string {
	return AsStrings(list, YM)
}
func AsYStrings(list []time.Time) []string {
	return AsStrings(list, Y)
}

// AsString converts a slice of times to a slice of strings
func AsStrings(set []time.Time, layout Format) (list []string) {
	var format = string(layout)
	list = []string{}
	for _, t := range set {
		list = append(list, t.Format(format))
	}
	return
}

// AsString convert time string to a date formatted string
func AsString(t time.Time, layout Format) string {
	return t.Format(string(layout))
}

// FromString converts string into a time where possible
func FromString(str string) (t time.Time, err error) {
	t, err = time.Parse(GetFormat(str), str)
	return
}
