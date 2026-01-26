package times

import "time"

type timeish struct {
	time.Time
}

func fromT(t time.Time) (result timeish, err error) {
	result = timeish{Time: t}
	return
}

func fromS(str string) (result timeish, err error) {
	t, err := time.Parse(GetFormat(str), string(FULL))
	// t, err := time.Parse(GetFormat(str), str) // ??
	if err == nil {
		result, err = fromT(t)
	}
	return
}
