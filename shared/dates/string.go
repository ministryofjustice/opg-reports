package dates

import "time"

func MustTime(t time.Time, e error) time.Time {
	if e != nil {
		return ErrorTime
	}
	return t
}

func ToTime(s string) (t time.Time, err error) {
	layout := GetFormat(s)
	t, err = time.Parse(layout, s)
	return
}

func Time(s string) time.Time {
	return MustTime(ToTime(s))
}
