package types

import "opg-reports/report/packages/types/interfaces"

var (
	_ interfaces.Insertable = &Team{}
	_ interfaces.Selectable = &Team{}
	_ interfaces.Statement  = &Select{}
)
