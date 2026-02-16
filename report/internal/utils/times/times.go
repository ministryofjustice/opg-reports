package times

import "time"

func AsTime(str string) (t time.Time, err error) {
	t, err = time.Parse(GetFormat(str), str)
	return
}
