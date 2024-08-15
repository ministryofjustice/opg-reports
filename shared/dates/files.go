package dates

import (
	"time"

	"github.com/djherbis/times"
)

func FileCreationTime(filepath string) (c time.Time, err error) {
	var t times.Timespec
	if t, err = times.Stat(filepath); err == nil && t.HasBirthTime() {
		c = t.BirthTime()
	}
	return

}
