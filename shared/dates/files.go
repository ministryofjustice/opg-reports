package dates

import (
	"os"
	"time"
)

func FileCreationTime(filepath string) (c time.Time, err error) {

	info, err := os.Stat(filepath)
	if err == nil {
		c = info.ModTime()
	}
	return

}
