package times

import (
	"fmt"
	"strings"
	"time"
)

// AsYMDString returns a YYYY-MM-DD formatted string
func AsYMDString(t time.Time) string {
	return AsString(t, YMD)
}

// AsYMString returns a YYYY-MM formatted string
func AsYMString(t time.Time) string {
	return AsString(t, YM)
}

// AsYString returns a YYYY formatted string
func AsYString(t time.Time) string {
	return AsString(t, Y)
}

// AsYMDStrings returns a slice of YYYY-MM-DD strings
func AsYMDStrings(list []time.Time) []string {
	return AsStrings(list, YMD)
}

// AsYMStrings returns a slice of YYYY-MM strings
func AsYMStrings(list []time.Time) []string {
	return AsStrings(list, YM)
}

// AsYStrings returns a slice of YYYY strings
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

// JoinedYMDList converts slice of times to a single string and list of strings
func JoinedYMList(months []time.Time) (str string, mths []string) {
	mths = AsYMStrings(months)
	str = fmt.Sprintf("'%s'", strings.Join(mths, "','"))
	return
}
