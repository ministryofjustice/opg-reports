package cliargs

import (
	"fmt"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/utils/convert"
)

type DateTime struct {
	Value *time.Time
}

func (self *DateTime) String() string {
	if self.Value != nil {
		return self.Value.Format(time.RFC3339)
	}
	return ""
}

func (self *DateTime) Set(s string) (err error) {
	var asTime time.Time

	asTime, err = convert.StringToTime(s)
	if err == nil {
		self.Value = &asTime
	} else {
		err = fmt.Errorf("failed to parse value [%s] into a time", s)
	}

	return
}
