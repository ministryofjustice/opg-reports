// Package dates provides helpful formatting constants and a convertor
//
// To keep date formats consistent between packages, we share the formats here and include
// a StringToDate func to easily convert a string into a time.Time
package dates

import "time"

const Format string = time.RFC3339
const FormatYMD string = "2006-01-02"
const FormatYM string = "2006-01"
const FormatY string = "2006"

func GetFormat(value string) string {
	max := len(Format)
	l := len(value)
	if l > max {
		return Format
	}

	return Format[:l]
}

func StringToDate(value string) (t time.Time, err error) {
	layout := GetFormat(value)
	t, err = time.Parse(layout, value)
	if err == nil {
		t = t.UTC()
	}
	return t, err
}

func StringToDateDefault(value string, comp string, defaultV string) (t time.Time, err error) {
	if value == comp {
		value = defaultV
	}
	return StringToDate(value)
}

func Strings(dates []time.Time, dateFormat string) []string {
	strs := []string{}
	for _, d := range dates {
		strs = append(strs, d.Format(dateFormat))
	}
	return strs
}
