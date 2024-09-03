package awsu

import (
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

func (a *AwsUptime) Insertable() InsertParams {
	bytes, _ := convert.Marshal(a)
	ip, _ := convert.Unmarshal[InsertParams](bytes)
	return ip
}
