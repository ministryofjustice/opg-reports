package awsc

import (
	"fmt"

	"github.com/ministryofjustice/opg-reports/shared/convert"
)

func (a *AwsCost) UID() string {
	return fmt.Sprintf("%d", a.ID)
}

func (a *AwsCost) Insertable() InsertParams {
	bytes, _ := convert.Marshal(a)
	ip, _ := convert.Unmarshal[InsertParams](bytes)
	return ip
}
