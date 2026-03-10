package types

import (
	"opg-reports/report/packages/types/interfaces"
)

// NilRow is an empty model that adheres to selectable interface
// so nothing is returned byt interface is for api is met
type NilRow struct{}

func (self *NilRow) Sequence() []any {
	return []any{}
}
func (self *NilRow) Result() interfaces.Result {
	return map[string]interface{}{}
}
